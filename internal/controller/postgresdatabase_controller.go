package controller

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	databasesv1 "github.com/mycompany/postgres-database-controller/api/v1"
)

// PlatformConfig holds platform-wide defaults
type PlatformConfig struct {
	DefaultStorageClass    string
	DefaultImageRegistry   string
	DefaultCRVersion       string
	DefaultPGBouncerImage   string
	DefaultPGBackRestImage string
	DefaultPMMImage        string
	DefaultPMMHost         string
}

// NewDefaultPlatformConfig returns platform defaults
func NewDefaultPlatformConfig() *PlatformConfig {
	return &PlatformConfig{
		DefaultStorageClass:    "standard",
		DefaultImageRegistry:   "docker.io/percona",
		DefaultCRVersion:       "2.8.2",
		DefaultPGBouncerImage:  "percona-pgbouncer:1.25.0-1",
		DefaultPGBackRestImage: "percona-pgbackrest:2.57.0-1",
		DefaultPMMImage:        "pmm-client:3.5.0",
		DefaultPMMHost:         "prometheus.monitoring",
	}
}

// PostgresDatabaseReconciler reconciles a PostgresDatabase object
type PostgresDatabaseReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Log           logr.Logger
	PlatformConfig *PlatformConfig
}

//+kubebuilder:rbac:groups=databases.mycompany.com,resources=postgresdatabases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=databases.mycompany.com,resources=postgresdatabases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=databases.mycompany.com,resources=postgresdatabases/finalizers,verbs=update
//+kubebuilder:rbac:groups=pgv2.percona.com,resources=perconapgclusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pgv2.percona.com,resources=perconapgclusters/status,verbs=get

// Reconcile implements the reconciliation logic for PostgresDatabase
func (r *PostgresDatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("postgresdatabase", req.NamespacedName)

	// Fetch the PostgresDatabase instance
	var db databasesv1.PostgresDatabase
	if err := r.Get(ctx, req.NamespacedName, &db); err != nil {
		if errors.IsNotFound(err) {
			log.Info("PostgresDatabase resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get PostgresDatabase")
		return ctrl.Result{}, err
	}

	// Generate the PerconaPGCluster from PostgresDatabase
	cluster := r.generatePerconaPGCluster(&db)

	// For now, just log what we would create
	log.Info("Would create PerconaPGCluster", "name", cluster.Name, "namespace", cluster.Namespace)
	
	// Update status to show we're working
	r.updateStatus(&db, "Creating", "PerconaPGCluster generation started", cluster.Name)

	return ctrl.Result{RequeueAfter: 30}, nil
}

// generatePerconaPGCluster creates a PerconaPGCluster from PostgresDatabase
func (r *PostgresDatabaseReconciler) generatePerconaPGCluster(db *databasesv1.PostgresDatabase) map[string]interface{} {
	config := r.PlatformConfig
	
	// Create the PerconaPGCluster manifest as a generic map
	cluster := map[string]interface{}{
		"apiVersion": "pgv2.percona.com/v2",
		"kind":       "PerconaPGCluster",
		"metadata": map[string]interface{}{
			"name":      db.Name,
			"namespace": db.Namespace,
			"labels": map[string]string{
				"created-by": "postgres-database-controller",
				"app":        "postgres-database",
			},
		},
		"spec": map[string]interface{}{
			"crVersion":     config.DefaultCRVersion,
			"image":         fmt.Sprintf("%s/percona-distribution-postgresql:%d.7-2", config.DefaultImageRegistry, db.Spec.Version),
			"postgresVersion": db.Spec.Version,
			"instances": []map[string]interface{}{
				{
					"name":     "instance1",
					"replicas": db.Spec.Replicas,
					"dataVolumeClaimSpec": map[string]interface{}{
						"accessModes": []string{"ReadWriteOnce"},
						"resources": map[string]interface{}{
							"requests": map[string]interface{}{
								"storage": db.Spec.Storage,
							},
						},
						"storageClassName": config.DefaultStorageClass,
					},
					"affinity": r.generatePodAntiAffinity(),
				},
			},
		},
	}

	// Add PgBouncer proxy
	cluster["spec"].(map[string]interface{})["proxy"] = map[string]interface{}{
		"pgBouncer": map[string]interface{}{
			"replicas": db.Spec.Replicas,
			"image":    fmt.Sprintf("%s/%s", config.DefaultImageRegistry, config.DefaultPGBouncerImage),
			"affinity": r.generatePodAntiAffinity(),
		},
	}

	// Add backups if enabled
	if db.Spec.Backup {
		cluster["spec"].(map[string]interface{})["backups"] = map[string]interface{}{
			"pgbackrest": map[string]interface{}{
				"image": fmt.Sprintf("%s/%s", config.DefaultImageRegistry, config.DefaultPGBackRestImage),
				"repos": []map[string]interface{}{{"name": "repo1"}},
			},
		}
	}

	// Add monitoring if enabled
	if db.Spec.Monitoring {
		cluster["spec"].(map[string]interface{})["pmm"] = map[string]interface{}{
			"enabled":    true,
			"image":      fmt.Sprintf("%s/%s", config.DefaultImageRegistry, config.DefaultPMMImage),
			"serverHost": config.DefaultPMMHost,
		},
	}

	return cluster
}

// generatePodAntiAffinity creates default pod anti-affinity rules
func (r *PostgresDatabaseReconciler) generatePodAntiAffinity() map[string]interface{} {
	return map[string]interface{}{
		"podAntiAffinity": map[string]interface{}{
			"preferredDuringSchedulingIgnoredDuringExecution": []map[string]interface{}{
				{
					"weight": 1,
					"podAffinityTerm": map[string]interface{}{
						"labelSelector": map[string]interface{}{
							"matchLabels": map[string]string{
								"postgres-operator.crunchydata.com/data": "postgres",
							},
						},
						"topologyKey": "kubernetes.io/hostname",
					},
				},
			},
		},
	}
}

// updateStatus updates the PostgresDatabase status
func (r *PostgresDatabaseReconciler) updateStatus(db *databasesv1.PostgresDatabase, phase, message, endpoint string) {
	db.Status.Phase = phase
	db.Status.Message = message
	db.Status.Endpoint = endpoint
	
	if err := r.Status().Update(context.Background(), db); err != nil {
		r.Log.Error(err, "Failed to update PostgresDatabase status")
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PostgresDatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasesv1.PostgresDatabase{}).
		Complete(r)
}