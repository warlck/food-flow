{{- define "food-flow.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "food-flow.fullname" -}}
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

{{- define "food-flow.componentName" -}}
{{- printf "%s-%s" (include "food-flow.fullname" .root) .name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "food-flow.serviceName" -}}
{{- if .overrideName }}
{{- .overrideName | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- include "food-flow.componentName" (dict "root" .root "name" .name) }}
{{- end }}
{{- end }}

{{- define "food-flow.configName" -}}
{{- printf "%s-%s" (include "food-flow.fullname" .root) .name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "food-flow.databaseSecretName" -}}
{{- if .Values.database.auth.existingSecret }}
{{- .Values.database.auth.existingSecret }}
{{- else }}
{{- include "food-flow.componentName" (dict "root" . "name" "database-credentials") }}
{{- end }}
{{- end }}

{{- define "food-flow.databaseClaimName" -}}
{{- if .Values.database.persistence.existingClaim }}
{{- .Values.database.persistence.existingClaim }}
{{- else }}
{{- include "food-flow.componentName" (dict "root" . "name" "database-data") }}
{{- end }}
{{- end }}

{{- define "food-flow.stripeSecretName" -}}
{{- if .Values.sales.stripe.secretName }}
{{- .Values.sales.stripe.secretName }}
{{- else }}
{{- include "food-flow.componentName" (dict "root" . "name" "stripe-secrets") }}
{{- end }}
{{- end }}

{{- define "food-flow.labels" -}}
app.kubernetes.io/name: {{ include "food-flow.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" }}
{{- end }}
