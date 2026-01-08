package v1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PostgresDatabase is the Schema for the postgresdatabases API
type PostgresDatabase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PostgresDatabaseSpec   `json:"spec,omitempty"`
	Status PostgresDatabaseStatus `json:"status,omitempty"`
}

// PostgresDatabaseSpec defines the desired state of PostgresDatabase
type PostgresDatabaseSpec struct {
	// PostgreSQL major version (13, 14, 15, 16, or 17)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=13;14;15;16;17
	Version int `json:"version"`

	// Number of database replicas (1 = single, 3+ = HA)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	Replicas int `json:"replicas"`

	// Storage capacity (e.g., 100Gi, 1Ti)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^[0-9]+[KMGT]i$
	Storage string `json:"storage"`

	// Enable automated backups to S3
	// +kubebuilder:default=true
	Backup bool `json:"backup,omitempty"`

	// Enable PMM monitoring
	// +kubebuilder:default=true
	Monitoring bool `json:"monitoring,omitempty"`

	// Resource limits and requests
	Resources *ResourceRequirements `json:"resources,omitempty"`

	// Backup retention period (e.g., 7d, 30d)
	// +kubebuilder:validation:Pattern=^[0-9]+[dwmy]$
	// +kubebuilder:default="7d"
	BackupRetention string `json:"backupRetention,omitempty"`
}

// ResourceRequirements defines CPU and memory requirements
type ResourceRequirements struct {
	// Resource requests
	Requests *ResourceList `json:"requests,omitempty"`
	// Resource limits
	Limits *ResourceList `json:"limits,omitempty"`
}

// ResourceList defines resource quantities
type ResourceList struct {
	// CPU requirement (e.g., 100m, 1)
	// +kubebuilder:validation:Pattern=^[0-9]+m?$
	CPU string `json:"cpu,omitempty"`
	
	// Memory requirement (e.g., 256Mi, 1Gi)
	// +kubebuilder:validation:Pattern=^[0-9]+[KMGT]i$
	Memory string `json:"memory,omitempty"`
}

// PostgresDatabaseStatus defines the observed state of PostgresDatabase
type PostgresDatabaseStatus struct {
	// Current phase of the database (Pending, Creating, Ready, Failed)
	// +kubebuilder:validation:Enum=Pending;Creating;Ready;Failed
	Phase string `json:"phase,omitempty"`

	// Connection endpoint for the database
	Endpoint string `json:"endpoint,omitempty"`

	// Actual number of ready replicas
	Replicas int32 `json:"replicas,omitempty"`

	// Status message or error details
	Message string `json:"message,omitempty"`

	// Name of secret containing connection credentials
	CredentialsSecret string `json:"credentialsSecret,omitempty"`

	// Timestamp of last successful backup
	LastBackupTime *metav1.Time `json:"lastBackupTime,omitempty"`

	// PostgreSQL version actually running
	RunningVersion string `json:"runningVersion,omitempty"`

	// Cluster state from underlying operator
	ClusterState string `json:"clusterState,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PostgresDatabaseList contains a list of PostgresDatabase
type PostgresDatabaseList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ListMeta   `json:"metadata,omitempty"`
	Items           []PostgresDatabase `json:"items"`
}

// SetDefaults sets default values for the PostgresDatabase spec
func (db *PostgresDatabase) SetDefaults() {
	if db.Spec.BackupRetention == "" {
		db.Spec.BackupRetention = "7d"
	}
	
	// Set default values if not specified
	if db.Spec.Resources == nil {
		db.Spec.Resources = &ResourceRequirements{
			Requests: &ResourceList{
				CPU:    "100m",
				Memory: "256Mi",
			},
			Limits: &ResourceList{
				CPU:    "500m",
				Memory: "1Gi",
			},
		}
	}
}

// Validate validates the PostgresDatabase spec
func (db *PostgresDatabase) Validate() error {
	// Validate storage quantity
	if _, err := resource.ParseQuantity(db.Spec.Storage); err != nil {
		return fmt.Errorf("invalid storage quantity: %v", err)
	}
	
	// Validate resource quantities if specified
	if db.Spec.Resources != nil {
		if db.Spec.Resources.Requests != nil {
			if db.Spec.Resources.Requests.CPU != "" {
				if _, err := resource.ParseQuantity(db.Spec.Resources.Requests.CPU); err != nil {
					return fmt.Errorf("invalid CPU request: %v", err)
				}
			}
			if db.Spec.Resources.Requests.Memory != "" {
				if _, err := resource.ParseQuantity(db.Spec.Resources.Requests.Memory); err != nil {
					return fmt.Errorf("invalid memory request: %v", err)
				}
			}
		}
		
		if db.Spec.Resources.Limits != nil {
			if db.Spec.Resources.Limits.CPU != "" {
				if _, err := resource.ParseQuantity(db.Spec.Resources.Limits.CPU); err != nil {
					return fmt.Errorf("invalid CPU limit: %v", err)
				}
			}
			if db.Spec.Resources.Limits.Memory != "" {
				if _, err := resource.ParseQuantity(db.Spec.Resources.Limits.Memory); err != nil {
					return fmt.Errorf("invalid memory limit: %v", err)
				}
			}
		}
	}
	
	return nil
}