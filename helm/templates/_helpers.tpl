{{- define "service.name" -}}
{{- .Chart.Name -}}
{{- end -}}

{{- define "service.fullname" -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "service.labels" -}}
helm.sh/chart: {{ include "service.name" . }}-{{ .Chart.Version | replace "+" "_" }}
{{- end -}}