{{/*
Expand the name of the chart.
*/}}
{{- define "gateway-extension.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "gateway-extension.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Namespace for all resources to be installed into
If not defined in values file then the helm release namespace is used
By default this is not set so the helm release namespace will be used

This gets around an problem within helm discussed here
https://github.com/helm/helm/issues/5358
*/}}
{{- define "gateway-extension.namespace" -}}
    {{ .Values.namespace | default .Release.Namespace }}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "gateway-extension.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "gateway-extension.labels" -}}
helm.sh/chart: {{ include "gateway-extension.chart" . }}
{{ include "gateway-extension.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "gateway-extension.selectorLabels" -}}
app.kubernetes.io/name: {{ include "gateway-extension.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: {{ .Chart.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "gateway-extension.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "gateway-extension.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Util function for generating the image URL based on the provided options.
*/}}
{{- define "gateway-extension.image" -}}
{{- $defaultTag := index . 1 -}}
{{- with index . 0 -}}
{{- if .registry -}}{{ printf "%s/%s" .registry .repository }}{{- else -}}{{- .repository -}}{{- end -}}
{{- if .digest -}}{{ printf "@%s" .digest }}{{- else -}}{{ printf ":%s" (default $defaultTag .tag) }}{{- end -}}
{{- end }}
{{- end }}

{{- define "gateway-extension.podsecuritycontext" }}
{{-   $vals := dict }}
{{-   if and (dig "recommendations" "enabled" false .Values.global) (dig "recommendations" "securityContexts" true .Values.global) }}
{{-       $vals = mergeOverwrite $vals (include "gateway-extension.default_k8s_podsecuritycontext_values" . | fromYaml) }}
{{-   end }}
{{-   $vals = mergeOverwrite $vals .Values.podSecurityContext }}
{{   toYaml $vals }}
{{- end }}

{{- define "gateway-extension.default_k8s_podsecuritycontext_values" }}
fsGroupChangePolicy: OnRootMismatch
runAsUser: 1000
runAsGroup: 1000
fsGroup: 1000
{{- end }}
