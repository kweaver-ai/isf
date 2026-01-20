{{/* vim: set filetype=mustache: */}}
{{- define "audit-log.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "audit-log.fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s" .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "audit-log.server" -}}
{{- printf "%s" (include "audit-log.fullname" .) -}}
{{- end -}}


{{- define "audit-log.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "audit-log.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}
