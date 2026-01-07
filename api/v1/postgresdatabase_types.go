package v1

import (
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
	// PostgreSQL major version
	Version int `json:"version"`

	// Number of database replicas (1 = single, 3 = HA)
	Replicas int `json:"replicas"`

	// Storage capacity (e.g., 100Gi, 1Ti)
	Storage string `json:"storage"`

	// Enable automated backups to S3
	Backup bool `json:"backup,omitempty"`

	// Enable PMM monitoring
	Monitoring bool `json:"monitoring,omitempty"`

	// Resource limits and requests
	Resources ResourceRequirements `json:"resources,omitempty"`

	// Backup retention period (e.g., 7d, 30d)
	BackupRetention string `json:"backupRetention,omitempty"`
}

// ResourceRequirements defines CPU and memory requirements
type ResourceRequirements struct {
	Requests ResourceList `json:"requests,omitempty"`
	Limits   ResourceList `json:"limits,omitempty"`
}

// ResourceList defines resource quantities
type ResourceList struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// PostgresDatabaseStatus defines the observed state of PostgresDatabase
type PostgresDatabaseStatus struct {
	// Current phase of the database
	Phase string `json:"phase,omitempty"`

	// Connection endpoint for the database
	Endpoint string `json:"endpoint,omitempty"`

	// Actual number of ready replicas
	Replicas int `json:"replicas,omitempty"`

	// Status message or error details
	Message string `json:"message,omitempty"`

	// Name of secret containing connection credentials
	CredentialsSecret string `json:"credentialsSecret,omitempty"`

	// Timestamp of last successful backup
	LastBackupTime *metav1.Time `json:"lastBackupTime,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PostgresDatabaseList contains a list of PostgresDatabase
type PostgresDatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgresDatabase `json:"items"`
}
