package postgresql

import (
	"os"
	"strconv"

	databasev1alpha1 "github.com/iboware/postgresql-operator/apis/database/v1alpha1"
	utils "github.com/iboware/postgresql-operator/utils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

//NewStatefulSetForCR creates a statefulSet of PostgreSQL
func NewStatefulSetForCR(cr *databasev1alpha1.PostgreSQL) *appsv1.StatefulSet {
	var log = logf.Log.WithName("controller_postgresql")

	labels := utils.LabelsForPostgreSQL(cr.ObjectMeta.Name)
	envVars := envVarsForPostgreSQL(cr.ObjectMeta.Name, int(cr.Spec.Replicas), cr.Spec.Namespace)
	livenessProbeCmd := []string{
		"sh",
		"-c",
		"PGPASSWORD=$POSTGRES_PASSWORD psql -w -U \"postgres\" -d \"postgres\"  -h 127.0.0.1 -c \"SELECT 1\"",
	}
	// Constants for StatefulSet & Volumes
	const (
		AppImage                     = "docker.io/bitnami/postgresql-repmgr:12-debian-10" //"nginxdemos/nginx-hello:latest"
		AppContainerName             = "postgresql"
		ImagePullPolicy              = v1.PullIfNotPresent
		CMDFileMode      os.FileMode = 0755
	)

	quantity, err := resource.ParseQuantity(cr.Spec.DiskSize)
	if err != nil {
		log.Error(err, "Invalid Disk Size %s", cr.Spec.DiskSize)
	}
	var (
		// storageClassName              = "standard"
		terminationGracePeriodSeconds = int64(10)
		fsGroup                       = int64(1001)
		volumeMode                    = int32(CMDFileMode)
		accessMode                    = []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce}
		resourceList                  = v1.ResourceList{v1.ResourceStorage: quantity}
	)

	statefulset := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.ObjectMeta.Name + "-postgresql",
			Namespace: cr.Spec.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: cr.ObjectMeta.Name + "-postgresql-headless",
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{Type: appsv1.RollingUpdateStatefulSetStrategyType},
			Replicas:       &cr.Spec.Replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
					SecurityContext:               &v1.PodSecurityContext{FSGroup: &fsGroup},
					Containers: []v1.Container{
						{
							Name:            AppContainerName,
							Image:           AppImage,
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
										Command: livenessProbeCmd,
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
								TimeoutSeconds:      5,
								SuccessThreshold:    1,
								FailureThreshold:    6,
							},
							ReadinessProbe: &v1.Probe{
								Handler: v1.Handler{
									Exec: &v1.ExecAction{
										Command: livenessProbeCmd,
									},
								},
								InitialDelaySeconds: 5,
								PeriodSeconds:       10,
								TimeoutSeconds:      5,
								SuccessThreshold:    1,
								FailureThreshold:    6,
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/bitnami/postgresql",
								},
								{
									Name:      "hooks-scripts",
									MountPath: "/pre-stop.sh",
									SubPath:   "pre-stop.sh",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "hooks-scripts",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									DefaultMode: &volumeMode,
									LocalObjectReference: v1.LocalObjectReference{
										Name: cr.ObjectMeta.Name + "-postgresql-hooks-scripts",
									},
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "data",
					},
					Spec: v1.PersistentVolumeClaimSpec{
						AccessModes: accessMode,
						Resources: v1.ResourceRequirements{
							Requests: resourceList,
						},
					},
				},
			},
		},
	}
	return statefulset
}

func envVarsForPostgreSQL(releaseName string, replicaCount int, namespace string) []v1.EnvVar {

	var partnerNodes string

	for i := 0; i < replicaCount; i++ {
		partnerNodes = partnerNodes + releaseName + "-postgresql-" + strconv.Itoa(i) + "." + releaseName + "-postgresql-headless." + namespace + ".svc.cluster.local,"
	}

	envVars := []v1.EnvVar{
		{
			Name:  "BITNAMI_DEBUG",
			Value: "false",
		},
		//PostgreSQL configuration
		{
			Name:  "POSTGRESQL_VOLUME_DIR",
			Value: "/bitnami/postgresql",
		},
		{
			Name:  "PGDATA",
			Value: "/bitnami/postgresql/data",
		},
		{
			Name:  "POSTGRES_USER",
			Value: "postgres",
		},
		{
			Name: "POSTGRES_PASSWORD",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: releaseName + "-postgresql"},
					Key:                  "postgresql-password",
				},
			},
		},
		{
			Name:  "POSTGRES_DB",
			Value: "postgres",
		},
		{
			Name: "MY_POD_NAME",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name:  "REPMGR_UPGRADE_EXTENSION",
			Value: "no",
		},
		{
			Name:  "REPMGR_PGHBA_TRUST_ALL",
			Value: "no",
		},
		{
			Name:  "REPMGR_MOUNTED_CONF_DIR",
			Value: "/bitnami/repmgr/conf",
		},
		{
			Name:  "REPMGR_PARTNER_NODES",
			Value: partnerNodes,
		},
		{
			Name:  "REPMGR_PRIMARY_HOST",
			Value: releaseName + "-postgresql-0." + releaseName + "-postgresql-headless." + namespace + ".svc.cluster.local",
		},
		{
			Name:  "REPMGR_NODE_NAME",
			Value: "$(MY_POD_NAME)",
		},
		{
			Name:  "REPMGR_NODE_NETWORK_NAME",
			Value: "$(MY_POD_NAME)." + releaseName + "-postgresql-headless." + namespace + ".svc.cluster.local",
		},
		{
			Name:  "REPMGR_LOG_LEVEL",
			Value: "NOTICE",
		},
		{
			Name:  "REPMGR_CONNECT_TIMEOUT",
			Value: "5",
		},
		{
			Name:  "REPMGR_RECONNECT_ATTEMPTS",
			Value: "3",
		},
		{
			Name:  "REPMGR_RECONNECT_INTERVAL",
			Value: "5",
		},
		{
			Name:  "REPMGR_USERNAME",
			Value: "repmgr",
		},
		{
			Name:  "REPMGR_DATABASE",
			Value: "repmgr"},
		{
			Name: "REPMGR_PASSWORD",
			ValueFrom: &v1.EnvVarSource{SecretKeyRef: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{Name: releaseName + "-postgresql"},
				Key:                  "repmgr-password",
			}},
		},
	}

	return envVars
}
