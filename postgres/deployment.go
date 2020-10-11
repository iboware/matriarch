package postgres

import (
	"os"
	"strconv"

	databasev1alpha1 "github.com/iboware/postgresql-operator/apis/database/v1alpha1"
	utils "github.com/iboware/postgresql-operator/utils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//NewPgPoolDeployment creates a deployment of PgPool
func NewPgPoolDeployment(cr *databasev1alpha1.PostgreSQL) *appsv1.Deployment {
	// var log = logf.Log.WithName("controller_pgpool")

	labels := utils.NewLabels("pgpool", cr.ObjectMeta.Name, "pgpool")
	envVars := EnvVarsForPgPool(cr.ObjectMeta.Name, int(cr.Spec.Replicas), cr.Spec.Namespace)
	readinessProbe := []string{
		"sh",
		"-ec",
		"PGPASSWORD=${PGPOOL_POSTGRES_PASSWORD} psql -U \"postgres\" -d \"postgres\" -h 127.0.0.1 -tA -c \"SELECT 1\" > /dev/null",
	}
	// Constants for StatefulSet & Volumes
	const (
		PostgresImage                     = "docker.io/bitnami/postgresql-repmgr:12-debian-10"
		PgPoolImage                       = "docker.io/bitnami/pgpool:4-debian-10"
		PostgresContainerName             = "postgresql"
		PgPoolContainerName               = "pgpool"
		ImagePullPolicy                   = v1.PullIfNotPresent
		CMDFileMode           os.FileMode = 0755
	)

	var (
		// storageClassName              = "standard"
		terminationGracePeriodSeconds = int64(10)
		fsGroup                       = int64(1001)
	)

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.ObjectMeta.Name + "-pgpool",
			Namespace: cr.Spec.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Replicas: &cr.Spec.Replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					SecurityContext:               &v1.PodSecurityContext{FSGroup: &fsGroup},
					Containers: []v1.Container{
						{
							Name:            PgPoolContainerName,
							Image:           PgPoolImage,
							ImagePullPolicy: ImagePullPolicy,
							SecurityContext: &v1.SecurityContext{RunAsUser: &fsGroup},
							Env:             envVars,
							Ports: []v1.ContainerPort{
								{Name: "postgresql",
									ContainerPort: 5432,
									Protocol:      v1.ProtocolTCP,
								},
							},
							LivenessProbe: &v1.Probe{
								Handler: v1.Handler{
									Exec: &v1.ExecAction{
										Command: []string{
											"/opt/bitnami/scripts/pgpool/healthcheck.sh",
										},
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
								TimeoutSeconds:      5,
								SuccessThreshold:    1,
								FailureThreshold:    5,
							},
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{
									Exec: &v1.ExecAction{
										Command: readinessProbe,
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       5,
								TimeoutSeconds:      5,
								SuccessThreshold:    1,
								FailureThreshold:    5,
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

// EnvVarsForPgPool Creates environmental Variables for PgPool
func EnvVarsForPgPool(releaseName string, replicaCount int, namespace string) []v1.EnvVar {

	var partnerNodes string

	for i := 0; i < replicaCount; i++ {
		partnerNodes = partnerNodes + strconv.Itoa(i) + ":" + releaseName + "-postgresql-" + strconv.Itoa(i) + "." + releaseName + "-postgresql-headless." + namespace + ".svc.cluster.local:5432,"
	}

	envVars := []v1.EnvVar{
		{
			Name:  "BITNAMI_DEBUG",
			Value: "false",
		},
		//PgPool configuration
		{
			Name:  "PGPOOL_ENABLE_LOAD_BALANCING",
			Value: "yes",
		},
		{
			Name:  "PGPOOL_ENABLE_LDAP",
			Value: "no",
		},
		{
			Name:  "PGPOOL_POSTGRES_USERNAME",
			Value: "postgres",
		},
		{
			Name: "PGPOOL_POSTGRES_PASSWORD",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: releaseName + "-postgresql"},
					Key:                  "postgresql-password",
				},
			},
		},
		{
			Name:  "PGPOOL_SR_CHECK_USER",
			Value: "repmgr",
		},
		{
			Name: "PGPOOL_SR_CHECK_PASSWORD",
			ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{Name: releaseName + "-postgresql"},
				Key:                  "repmgr-password",
			}},
		},
		{
			Name:  "PGPOOL_ADMIN_USERNAME",
			Value: "admin",
		},
		{
			Name: "PGPOOL_ADMIN_PASSWORD",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: releaseName + "-postgresql"},
					Key:                  "admin-password",
				},
			}},
		{
			Name:  "PGPOOL_BACKEND_NODES",
			Value: partnerNodes,
		},
	}

	return envVars
}
