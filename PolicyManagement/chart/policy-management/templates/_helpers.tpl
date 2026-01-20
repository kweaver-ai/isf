{{/* vim: set filetype=mustache: */}}

{{- define "policy-management.fullname" -}}
{{- $name := default .Chart.Name -}}
{{- printf "%s" .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "policy-management.server" -}}
{{- printf "%s" (include "policy-management.fullname" .) -}}
{{- end -}}
