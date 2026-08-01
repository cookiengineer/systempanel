package css

import "github.com/cookiengineer/systempanel/bindings/gtk4"

const Stylesheet = `
.sidebar-row {
	padding: 8px 12px;
	min-height: 36px;
}

.sidebar-row:disabled {
	opacity: 0.4;
}

.sidebar-row-label {
	font-size: 13px;
}

.volume-slider {
	margin: 4px 0;
}

.device-row {
	padding: 6px 8px;
}

.journal-emerg {
	color: #d64937;
	font-weight: bold;
}

.journal-alert {
	color: #d64937;
	font-weight: bold;
}

.journal-crit {
	color: #d64937;
	font-weight: bold;
}

.journal-err {
	color: #d64937;
}

.journal-warning {
	color: #e5a50a;
}

.journal-notice {
	color: #3584e4;
}

.journal-info {
	color: #33d17a;
}

.journal-debug {
	color: #888888;
}

.service-active {
	color: #33d17a;
	font-weight: bold;
}

.service-failed {
	color: #d64937;
	font-weight: bold;
}

.service-inactive {
	color: #888888;
}

.header-label {
	font-size: 18px;
	font-weight: bold;
}

.settings-row {
	padding: 4px 12px;
	margin: 2px 0;
}

.disk-health-good {
	color: #33d17a;
	font-weight: bold;
}

.disk-health-ok {
	color: #3584e4;
	font-weight: bold;
}

.disk-health-warning {
	color: #e5a50a;
	font-weight: bold;
}

.disk-health-critical {
	color: #d64937;
	font-weight: bold;
}

.group-header {
	font-size: 15px;
	font-weight: bold;
	margin-top: 8px;
}

.tag-pill {
	border-radius: 10px;
	padding: 2px 6px;
	margin: 2px;
	background: alpha(@theme_fg_color, 0.08);
}

.monospace-label {
	font-family: monospace;
}
`

func init() {
	provider := gtk4.CSSNew()
	provider.LoadFromString(Stylesheet)
	provider.ApplyToDisplay(gtk4.CSSPriorityApplication)
}
