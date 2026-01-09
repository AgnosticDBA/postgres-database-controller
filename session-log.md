🔍 Phase 1: Root Cause Analysis
Timestamp: Fri  9 Jan 2026 01:16:04 CET
Issues found:
Error: unknown flag: --dry-run
Usage:
  controller-gen [flags]

Examples:
	# Generate RBAC manifests and crds for all types under apis/,
	# outputting crds to /tmp/crds and everything else to stdout
	controller-gen rbac:roleName=<role name> crd paths=./apis/... output:crd:dir=/tmp/crds output:stdout

	# Generate deepcopy/runtime.Object implementations for a particular file
	controller-gen object paths=./apis/v1beta1/some_types.go

	# Generate OpenAPI v3 schemas for API packages and merge them into existing CRD manifests
	controller-gen schemapatch:manifests=./manifests output:dir=./manifests paths=./pkg/apis/... 

	# Run all the generators for a given project
	controller-gen paths=./apis/...

	# Explain the markers for generating CRDs, and their arguments
	controller-gen crd -ww


Flags:
  -h, --detailed-help count   print out more detailed help
                              (up to -hhh for the most detailed output, or -hhhh for json output)
      --help                  print out usage and a summary of options
      --version               show version
  -w, --which-markers count   print out all markers available with the requested generators
                              (up to -www for the most detailed output, or -wwww for json output)


Options


generators

+webhook[:headerFile=<string>][,year=<string>]                                                                                                                                                                                                        package  generates (partial) {Mutating,Validating}WebhookConfiguration objects.  
+schemapatch[:generateEmbeddedObjectMeta=<bool>],manifests=<string>[,maxDescLen=<int>]                                                                                                                                                                package  patches existing CRDs with new schemata.                                
+rbac[:headerFile=<string>],roleName=<string>[,year=<string>]                                                                                                                                                                                         package  generates ClusterRole objects.                                          
+object[:headerFile=<string>][,year=<string>]                                                                                                                                                                                                         package  generates code containing DeepCopy, DeepCopyInto, and                   
+crd[:allowDangerousTypes=<bool>][,crdVersions=<[]string>][,deprecatedV1beta1CompatibilityPreserveUnknownFields=<bool>][,generateEmbeddedObjectMeta=<bool>][,headerFile=<string>][,ignoreUnexportedFields=<bool>][,maxDescLen=<int>][,year=<string>]  package  generates CustomResourceDefinition objects.                             


generic

+paths=<[]string>  package  represents paths and go-style path patterns to use as package roots.  


output rules (optionally as output:<generator>:...)

+output:artifacts[:code=<string>],config=<string>  package  outputs artifacts to different locations, depending on    
+output:dir=<string>                               package  outputs each artifact to the given directory, regardless  
+output:none                                       package  skips outputting anything.                                
+output:stdout                                     package  outputs everything to standard-out, with no separation.   

run `controller-gen object paths=./api/v1/ --dry-run -w` to see all available markers, or `controller-gen object paths=./api/v1/ --dry-run -h` for usage
Build status: 0
Build status: 0
🧪 Phase 2: Core Fix
Timestamp: Fri  9 Jan 2026 01:16:43 CET
Error: unknown flag: --dry-run
Usage:
  controller-gen [flags]

Examples:
	# Generate RBAC manifests and crds for all types under apis/,
	# outputting crds to /tmp/crds and everything else to stdout
	controller-gen rbac:roleName=<role name> crd paths=./apis/... output:crd:dir=/tmp/crds output:stdout

	# Generate deepcopy/runtime.Object implementations for a particular file
	controller-gen object paths=./apis/v1beta1/some_types.go

	# Generate OpenAPI v3 schemas for API packages and merge them into existing CRD manifests
	controller-gen schemapatch:manifests=./manifests output:dir=./manifests paths=./pkg/apis/... 

	# Run all the generators for a given project
	controller-gen paths=./apis/...

	# Explain the markers for generating CRDs, and their arguments
	controller-gen crd -ww


Flags:
  -h, --detailed-help count   print out more detailed help
                              (up to -hhh for the most detailed output, or -hhhh for json output)
      --help                  print out usage and a summary of options
      --version               show version
  -w, --which-markers count   print out all markers available with the requested generators
                              (up to -www for the most detailed output, or -wwww for json output)


Options


generators

+webhook[:headerFile=<string>][,year=<string>]                                                                                                                                                                                                        package  generates (partial) {Mutating,Validating}WebhookConfiguration objects.  
+schemapatch[:generateEmbeddedObjectMeta=<bool>],manifests=<string>[,maxDescLen=<int>]                                                                                                                                                                package  patches existing CRDs with new schemata.                                
+rbac[:headerFile=<string>],roleName=<string>[,year=<string>]                                                                                                                                                                                         package  generates ClusterRole objects.                                          
+object[:headerFile=<string>][,year=<string>]                                                                                                                                                                                                         package  generates code containing DeepCopy, DeepCopyInto, and                   
+crd[:allowDangerousTypes=<bool>][,crdVersions=<[]string>][,deprecatedV1beta1CompatibilityPreserveUnknownFields=<bool>][,generateEmbeddedObjectMeta=<bool>][,headerFile=<string>][,ignoreUnexportedFields=<bool>][,maxDescLen=<int>][,year=<string>]  package  generates CustomResourceDefinition objects.                             


generic

+paths=<[]string>  package  represents paths and go-style path patterns to use as package roots.  


output rules (optionally as output:<generator>:...)

+output:artifacts[:code=<string>],config=<string>  package  outputs artifacts to different locations, depending on    
+output:dir=<string>                               package  outputs each artifact to the given directory, regardless  
+output:none                                       package  skips outputting anything.                                
+output:stdout                                     package  outputs everything to standard-out, with no separation.   

run `controller-gen object paths=./api/v1/ --dry-run -w` to see all available markers, or `controller-gen object paths=./api/v1/ --dry-run -h` for usage
📦 Phase 3: Integration Testing
Timestamp: Fri  9 Jan 2026 01:17:14 CET
🧪 Phase 2: Core Fix
Timestamp: Fri  9 Jan 2026 01:17:25 CET
