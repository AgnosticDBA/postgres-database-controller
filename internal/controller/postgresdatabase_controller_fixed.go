package controller

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	databasesv1 "github.com/mycompany/postgres-database-controller/api/v1"
)

// PostgresDatabaseReconciler reconciles a PostgresDatabase object
type PostgresDatabaseReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	Log            logr.Logger
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

	// Check if PerconaPGCluster already exists
	clusterName := types.NamespacedName{Name: db.Name, Namespace: db.Namespace}
	
	// Try to find existing cluster
	var existingCluster client.Object
	clusterGVK := metav1.GroupVersionKind{
		Group:   "pgv2.percona.com",
		Version: "v2",
		Kind:    "PerconaPGCluster",
	}
	
	existingCluster = &runtime.Unstructured{}
	existingCluster.GetObjectKind().SetGroupVersionKind(clusterGVK)
	
	if err := r.Get(ctx, clusterName, existingCluster); err != nil {
		if errors.IsNotFound(err) {
			// Create the PerconaPGCluster
			return r.createPerconaPGCluster(ctx, &db)
		}
		log.Error(err, "Failed to get existing PerconaPGCluster")
		return ctrl.Result{}, err
	}

	// Update status based on existing cluster
	return r.updateStatusFromCluster(ctx, &db, existingCluster)
}

// createPerconaPGCluster creates the underlying PerconaPGCluster
func (r *PostgresDatabaseReconciler) createPerconaPGCluster(ctx context.Context, db *databasesv1.PostgresDatabase) (ctrl.Result, error) {
	log := r.Log.WithValues("postgresdatabase", types.NamespacedName{Name: db.Name, Namespace: db.Namespace})

	// Generate the PerconaPGCluster YAML
	clusterYAML := r.generateClusterYAML(db)
	
	// Parse YAML into Unstructured object
	clusterObj := &runtime.Unstructured{}
	if err := yaml.Unmarshal([]byte(clusterYAML), clusterObj); err != nil {
		log.Error(err, "Failed to parse cluster YAML")
		return ctrl.Result{}, err
	}

	// Set the controller reference
	if err := ctrl.SetControllerReference(db, clusterObj, r.Scheme); err != nil {
		log.Error(err, "Failed to set controller reference")
		return ctrl.Result{}, err
	}

	// Create the PerconaPGCluster
	if err := r.Create(ctx, clusterObj); err != nil {
		log.Error(err, "Failed to create PerconaPGCluster")
		r.updateStatus(db, "Failed", fmt.Sprintf("Failed to create PerconaPGCluster: %v", err), "")
		return ctrl.Result{}, err
	}

	log.Info("Successfully created PerconaPGCluster", "name", db.Name)
	r.updateStatus(db, "Creating", "PerconaPGCluster created, waiting for ready status", fmt.Sprintf("%s-rw.%s.svc.cluster.local", db.Name, db.Namespace))

	return ctrl.Result{RequeueAfter: 30}, nil
}

// generateClusterYAML creates the PerconaPGCluster YAML from PostgresDatabase
func (r *PostgresDatabaseReconciler) generateClusterYAML(db *databasesv1.PostgresDatabase) string {
	config := r.PlatformConfig
	
	// Base cluster configuration
	clusterYAML := fmt.Sprintf(`
apiVersion: pgv2.percona.com/v2
kind: PerconaPGCluster
metadata:
  name: %s
  namespace: %s
  labels:
    created-by: postgres-database-controller
    app: postgres-database
spec:
  crVersion: "%s"
  image: %s/percona-distribution-postgresql:%d.7-2
  postgresVersion: %d
  instances:
  - name: instance1
    replicas: %d
    dataVolumeClaimSpec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: %s
      storageClassName: %s
    affinity:
      podAntiAffinity:
        preferredDuringSchedulingIgnoredDuringExecution:
        - weight: 1
          podAffinityTerm:
            labelSelector:
              matchLabels:
                postgres-operator.crunchydata.com/data: postgres
            topologyKey: kubernetes.io/hostname
  proxy:
    pgBouncer:
      replicas: %d
      image: %s/%s
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 1
            podAffinityTerm:
              labelSelector:
                matchLabels:
                  postgres-operator.crunchydata.com/role: pgbouncer
              topologyKey: kubernetes.io/hostname`,
		db.Name, db.Namespace, config.DefaultCRVersion,
		config.DefaultImageRegistry, db.Spec.Version, db.Spec.Version,
		db.Spec.Replicas, db.Spec.Storage, config.DefaultStorageClass,
		db.Spec.Replicas, config.DefaultImageRegistry, config.DefaultPGBouncerImage)

	// Add backups if enabled
	if db.Spec.Backup {
		clusterYAML += fmt.Sprintf(`
  backups:
    pgbackrest:
      image: %s/%s
      repos:
      - name: repo1
        schedules:
          full: "0 0 * * 6"`,
			config.DefaultImageRegistry, config.DefaultPGBackRestImage)
	}

	// Add monitoring if enabled
	if db.Spec.Monitoring {
		clusterYAML += fmt.Sprintf(`
  pmm:
    enabled: true
    image: %s/%s
    serverHost: %s`,
			config.DefaultImageRegistry, config.DefaultPMMImage, config.DefaultPMMHost)
	}

	return clusterYAML
}

// updateStatusFromCluster updates PostgresDatabase status based on PerconaPGCluster
func (r *PostgresDatabaseReconciler) updateStatusFromCluster(ctx context.Context, db *databasesv1.PostgresDatabase, cluster client.Object) (ctrl.Result, error) {
	log := r.Log.WithValues("postgresdatabase", types.NamespacedName{Name: db.Name, Namespace: db.Namespace})

	// Determine cluster status
	phase := "Creating"
	message := "PerconaPGCluster is being provisioned"
	endpoint := fmt.Sprintf("%s-rw.%s.svc.cluster.local", db.Name, db.Namespace)
	replicas := int32(0)

	// Check cluster status (simplified - in real implementation, check cluster.Status)
	annotations := cluster.GetAnnotations()
	if annotations != nil && annotations["postgres-operator.crunchydata.com/state"] == "Ready" {
		phase = "Ready"
		message = "PostgreSQL database is ready for connections"
		replicas = db.Spec.Replicas
	}

	// Update status
	db.Status.Replicas = replicas
	db.Status.Endpoint = endpoint
	db.Status.CredentialsSecret = fmt.Sprintf("%s.postgres-secret", db.Name)

	r.updateStatus(db, phase, message, endpoint)

	return ctrl.Result{RequeueAfter: 30}, nil
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