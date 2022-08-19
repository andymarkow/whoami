{{/*
	Expand the name of the chart.
*/}}
{{- define "whoami.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
	Create a default fully qualified app name.
	We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
	If release name contains chart name it will be used as a full name.
*/}}
{{- define "whoami.fullname" -}}
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
	Create chart name and version as used by the chart label.
*/}}
{{- define "whoami.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
	Common labels
*/}}
{{- define "whoami.labels" -}}
helm.sh/chart: {{ include "whoami.chart" . | quote }}
{{ include "whoami.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
{{- end }}

{{/*
	Selector labels
*/}}
{{- define "whoami.selectorLabels" -}}
app.kubernetes.io/name: {{ include "whoami.name" . | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
{{- end }}

{{/*
	Create the name of the service account to use
*/}}
{{- define "whoami.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "whoami.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
	Sets the application affinity for pod placement
*/}}
{{- define "whoami.affinity" -}}
  {{- if .Values.affinity }}
      affinity:
        {{ $tp := typeOf .Values.affinity }}
        {{- if eq $tp "string" }}
          {{- tpl .Values.affinity . | nindent 8 | trim }}
        {{- else }}
          {{- toYaml .Values.affinity | nindent 8 }}
        {{- end }}
  {{- end }}
{{- end -}}
