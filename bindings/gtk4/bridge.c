#include <gtk/gtk.h>
#include <stdint.h>

extern void goBridgeVoid(uintptr_t data);
extern void goBridgeListBoxRowActivated(uintptr_t data, void *row);
extern int goBridgeBool3Uint(uintptr_t data, unsigned int a, unsigned int b, unsigned int c);
extern int goBridgeBool(uintptr_t data);
extern void goBridgeDouble(uintptr_t data, double val);

static void _bridge_void(GtkWidget *widget, gpointer data) {
	goBridgeVoid((uintptr_t)data);
}

static void _bridge_app_activate(GtkApplication *app, gpointer data) {
	goBridgeVoid((uintptr_t)data);
}

static void _bridge_listbox_row_activated(GtkListBox *box, GtkListBoxRow *row, gpointer data) {
	goBridgeListBoxRowActivated((uintptr_t)data, (void *)row);
}

static gboolean _bridge_key_pressed(GtkEventControllerKey *ctrl, unsigned int keyval,
		unsigned int keycode, GdkModifierType state, gpointer data) {
	return goBridgeBool3Uint((uintptr_t)data, keyval, keycode, (unsigned int)state);
}

static gboolean _bridge_close_request(GtkWindow *win, gpointer data) {
	return goBridgeBool((uintptr_t)data);
}

static void _bridge_scale_value_changed(GtkRange *range, gpointer data) {
	goBridgeDouble((uintptr_t)data, gtk_range_get_value(range));
}

static void _bridge_notify(GObject *obj, GParamSpec *pspec, gpointer data) {
	goBridgeVoid((uintptr_t)data);
}

void connectGSignal(void *instance, const char *signal, int handlerType, uintptr_t handle) {
	GCallback cb = NULL;
	switch (handlerType) {
	case 0: cb = G_CALLBACK(_bridge_void); break;
	case 1: cb = G_CALLBACK(_bridge_app_activate); break;
	case 2: cb = G_CALLBACK(_bridge_listbox_row_activated); break;
	case 3: cb = G_CALLBACK(_bridge_key_pressed); break;
	case 4: cb = G_CALLBACK(_bridge_close_request); break;
	case 5: cb = G_CALLBACK(_bridge_scale_value_changed); break;
	case 6: cb = G_CALLBACK(_bridge_notify); break;
	}
	if (cb) {
		g_signal_connect_data(instance, signal, cb, (gpointer)(uintptr_t)handle, NULL, 0);
	}
}

void *gtk4Init(void) { gtk_init(); return NULL; }

void *gtk4AppNew(const char *id) { return gtk_application_new(id, G_APPLICATION_DEFAULT_FLAGS); }
int gtk4AppRun(void *app, int argc, char **argv) { return g_application_run(G_APPLICATION(app), argc, argv); }
void gtk4AppQuit(void *app) { g_application_quit(G_APPLICATION(app)); }
void gtk4AppAddWindow(void *app, void *win) { gtk_application_add_window((GtkApplication*)app, (GtkWindow*)win); }

void *gtk4WindowNew(void) { return gtk_window_new(); }
void *gtk4AppWindowNew(void *app) { return gtk_application_window_new((GtkApplication*)app); }
void gtk4WindowSetTitle(void *win, const char *title) { gtk_window_set_title((GtkWindow*)win, title); }
void gtk4WindowSetDefaultSize(void *win, int w, int h) { gtk_window_set_default_size((GtkWindow*)win, w, h); }
void gtk4WindowSetChild(void *win, void *child) { gtk_window_set_child((GtkWindow*)win, (GtkWidget*)child); }
void gtk4WindowSetTitlebar(void *win, void *bar) { gtk_window_set_titlebar((GtkWindow*)win, (GtkWidget*)bar); }
void gtk4WindowPresent(void *win) { gtk_window_present((GtkWindow*)win); }
void gtk4WindowClose(void *win) { gtk_window_close((GtkWindow*)win); }
void gtk4WindowDestroy(void *win) { gtk_window_destroy((GtkWindow*)win); }

void *gtk4HeaderBarNew(void) { return gtk_header_bar_new(); }
void gtk4HeaderBarSetTitleWidget(void *bar, void *w) { gtk_header_bar_set_title_widget((GtkHeaderBar*)bar, (GtkWidget*)w); }
void gtk4HeaderBarPackStart(void *bar, void *w) { gtk_header_bar_pack_start((GtkHeaderBar*)bar, (GtkWidget*)w); }
void gtk4HeaderBarPackEnd(void *bar, void *w) { gtk_header_bar_pack_end((GtkHeaderBar*)bar, (GtkWidget*)w); }
void gtk4HeaderBarSetShowTitleButtons(void *bar, int show) { gtk_header_bar_set_show_title_buttons((GtkHeaderBar*)bar, show); }

void *gtk4BoxNew(int orientation, int spacing) { return gtk_box_new(orientation, spacing); }
void gtk4BoxAppend(void *box, void *child) { gtk_box_append((GtkBox*)box, (GtkWidget*)child); }
void gtk4BoxPrepend(void *box, void *child) { gtk_box_prepend((GtkBox*)box, (GtkWidget*)child); }
void gtk4BoxRemove(void *box, void *child) { gtk_box_remove((GtkBox*)box, (GtkWidget*)child); }
void gtk4BoxSetSpacing(void *box, int spacing) { gtk_box_set_spacing((GtkBox*)box, spacing); }
void gtk4BoxSetHomogeneous(void *box, int v) { gtk_box_set_homogeneous((GtkBox*)box, v); }

void *gtk4ButtonNew(void) { return gtk_button_new(); }
void *gtk4ButtonNewWithLabel(const char *label) { return gtk_button_new_with_label(label); }
void gtk4ButtonSetLabel(void *btn, const char *label) { gtk_button_set_label((GtkButton*)btn, label); }
void gtk4ButtonSetIconName(void *btn, const char *icon) { gtk_button_set_icon_name((GtkButton*)btn, icon); }
void gtk4ButtonSetChild(void *btn, void *child) { gtk_button_set_child((GtkButton*)btn, (GtkWidget*)child); }

void *gtk4LabelNew(const char *text) { return gtk_label_new(text); }
void gtk4LabelSetText(void *lbl, const char *text) { gtk_label_set_text((GtkLabel*)lbl, text); }
void gtk4LabelSetMarkup(void *lbl, const char *markup) { gtk_label_set_markup((GtkLabel*)lbl, markup); }
void gtk4LabelSetWrap(void *lbl, int v) { gtk_label_set_wrap((GtkLabel*)lbl, v); }
void gtk4LabelSetXAlign(void *lbl, float a) { gtk_label_set_xalign((GtkLabel*)lbl, a); }
const char *gtk4LabelGetText(void *lbl) { return gtk_label_get_text((GtkLabel*)lbl); }

void *gtk4ImageNewFromIconName(const char *name) { return gtk_image_new_from_icon_name(name); }
void gtk4ImageSetIconName(void *img, const char *name) { gtk_image_set_from_icon_name((GtkImage*)img, name); }
void gtk4ImageSetPixelSize(void *img, int size) { gtk_image_set_pixel_size((GtkImage*)img, size); }

void *gtk4StackNew(void) { return gtk_stack_new(); }
void *gtk4StackAddTitled(void *s, void *child, const char *name, const char *title) {
	return gtk_stack_add_titled((GtkStack*)s, (GtkWidget*)child, name, title);
}
void gtk4StackSetVisibleChildName(void *s, const char *name) { gtk_stack_set_visible_child_name((GtkStack*)s, name); }
void gtk4StackRemove(void *s, void *child) { gtk_stack_remove((GtkStack*)s, (GtkWidget*)child); }
const char *gtk4StackGetVisibleChildName(void *s) { return gtk_stack_get_visible_child_name((GtkStack*)s); }
void *gtk4StackGetChildByName(void *s, const char *name) { return gtk_stack_get_child_by_name((GtkStack*)s, name); }
void gtk4StackSetTransitionType(void *s, int t) { gtk_stack_set_transition_type((GtkStack*)s, t); }
void gtk4StackSetTransitionDuration(void *s, unsigned int d) { gtk_stack_set_transition_duration((GtkStack*)s, d); }
void gtk4StackSetVHomogeneous(void *s, int v) { gtk_stack_set_vhomogeneous((GtkStack*)s, v); }
void *gtk4StackGetPage(void *s, void *child) { return gtk_stack_get_page((GtkStack*)s, (GtkWidget*)child); }
void gtk4StackPageSetTitle(void *page, const char *title) { gtk_stack_page_set_title((GtkStackPage*)page, title); }
void gtk4StackPageSetIconName(void *page, const char *name) { gtk_stack_page_set_icon_name((GtkStackPage*)page, name); }
void gtk4StackPageSetVisible(void *page, int v) { gtk_stack_page_set_visible((GtkStackPage*)page, v); }

void *gtk4ListBoxNew(void) { return gtk_list_box_new(); }
void gtk4ListBoxAppend(void *lb, void *row) { gtk_list_box_append((GtkListBox*)lb, (GtkWidget*)row); }
void gtk4ListBoxRemove(void *lb, void *row) { gtk_list_box_remove((GtkListBox*)lb, (GtkWidget*)row); }
void gtk4ListBoxSelectRow(void *lb, void *row) { gtk_list_box_select_row((GtkListBox*)lb, (GtkListBoxRow*)row); }
void *gtk4ListBoxGetSelectedRow(void *lb) { return gtk_list_box_get_selected_row((GtkListBox*)lb); }
void gtk4ListBoxSetSelectionMode(void *lb, int mode) { gtk_list_box_set_selection_mode((GtkListBox*)lb, mode); }

void *gtk4ListBoxRowNew(void) { return gtk_list_box_row_new(); }
void gtk4ListBoxRowSetChild(void *row, void *child) { gtk_list_box_row_set_child((GtkListBoxRow*)row, (GtkWidget*)child); }
int gtk4ListBoxRowGetIndex(void *row) { return gtk_list_box_row_get_index((GtkListBoxRow*)row); }

void *gtk4ScrolledWindowNew(void) { return gtk_scrolled_window_new(); }
void gtk4ScrolledWindowSetChild(void *sw, void *child) { gtk_scrolled_window_set_child((GtkScrolledWindow*)sw, (GtkWidget*)child); }
void gtk4ScrolledWindowSetPolicy(void *sw, int h, int v) { gtk_scrolled_window_set_policy((GtkScrolledWindow*)sw, h, v); }

void *gtk4ScaleNew(int orientation) {
	GtkAdjustment *adj = gtk_adjustment_new(0, 0, 100, 1, 10, 0);
	GtkWidget *w = gtk_scale_new(orientation, adj);
	gtk_scale_set_draw_value((GtkScale*)w, FALSE);
	return w;
}
void *gtk4ScaleNewWithRange(int orientation, double min, double max, double step) {
	GtkAdjustment *adj = gtk_adjustment_new(min, min, max, step, step*10, 0);
	GtkWidget *w = gtk_scale_new(orientation, adj);
	gtk_scale_set_draw_value((GtkScale*)w, FALSE);
	return w;
}
double gtk4ScaleGetValue(void *s) { return gtk_range_get_value((GtkRange*)s); }
void gtk4ScaleSetValue(void *s, double v) { gtk_range_set_value((GtkRange*)s, v); }
void gtk4ScaleSetRange(void *s, double min, double max) {
	GtkAdjustment *adj = gtk_range_get_adjustment((GtkRange*)s);
	gtk_adjustment_set_lower(adj, min);
	gtk_adjustment_set_upper(adj, max);
}

void *gtk4SwitchNew(void) { return gtk_switch_new(); }
void gtk4SwitchSetActive(void *sw, int v) { gtk_switch_set_active((GtkSwitch*)sw, v); }
int gtk4SwitchGetActive(void *sw) { return gtk_switch_get_active((GtkSwitch*)sw); }

void *gtk4EntryNew(void) { return gtk_entry_new(); }
void gtk4EntrySetText(void *e, const char *text) { gtk_entry_buffer_set_text(gtk_entry_get_buffer((GtkEntry*)e), text, -1); }
const char *gtk4EntryGetText(void *e) { return gtk_entry_buffer_get_text(gtk_entry_get_buffer((GtkEntry*)e)); }
void gtk4EntrySetPlaceholder(void *e, const char *text) { gtk_entry_set_placeholder_text((GtkEntry*)e, text); }

void *gtk4CssProviderNew(void) { return gtk_css_provider_new(); }
void gtk4CssLoadFromString(void *css, const char *data) { gtk_css_provider_load_from_string((GtkCssProvider*)css, data); }
void gtk4CssApplyToDisplay(void *css, unsigned int priority) {
	gtk_style_context_add_provider_for_display(gdk_display_get_default(), (GtkStyleProvider*)css, priority);
}

void *gtk4SpinnerNew(void) { return gtk_spinner_new(); }
void gtk4SpinnerStart(void *s) { gtk_spinner_start((GtkSpinner*)s); }
void gtk4SpinnerStop(void *s) { gtk_spinner_stop((GtkSpinner*)s); }

void *gtk4LevelBarNew(void) { return gtk_level_bar_new(); }
void gtk4LevelBarSetValue(void *lb, double v) { gtk_level_bar_set_value((GtkLevelBar*)lb, v); }

void *gtk4ControllerKeyNew(void) { return gtk_event_controller_key_new(); }
void gtk4WidgetAddController(void *w, void *ctrl) { gtk_widget_add_controller((GtkWidget*)w, (GtkEventController*)ctrl); }

void gtk4WidgetShow(void *w) { gtk_widget_set_visible((GtkWidget*)w, TRUE); }
void gtk4WidgetHide(void *w) { gtk_widget_set_visible((GtkWidget*)w, FALSE); }
void gtk4WidgetSetSensitive(void *w, int v) { gtk_widget_set_sensitive((GtkWidget*)w, v); }
void gtk4WidgetSetSizeRequest(void *w, int width, int height) { gtk_widget_set_size_request((GtkWidget*)w, width, height); }
void gtk4WidgetSetHExpand(void *w, int v) { gtk_widget_set_hexpand((GtkWidget*)w, v); }
void gtk4WidgetSetVExpand(void *w, int v) { gtk_widget_set_vexpand((GtkWidget*)w, v); }
void gtk4WidgetSetHAlign(void *w, int a) { gtk_widget_set_halign((GtkWidget*)w, a); }
void gtk4WidgetSetVAlign(void *w, int a) { gtk_widget_set_valign((GtkWidget*)w, a); }
void gtk4WidgetSetMarginStart(void *w, int m) { gtk_widget_set_margin_start((GtkWidget*)w, m); }
void gtk4WidgetSetMarginEnd(void *w, int m) { gtk_widget_set_margin_end((GtkWidget*)w, m); }
void gtk4WidgetSetMarginTop(void *w, int m) { gtk_widget_set_margin_top((GtkWidget*)w, m); }
void gtk4WidgetSetMarginBottom(void *w, int m) { gtk_widget_set_margin_bottom((GtkWidget*)w, m); }
void gtk4WidgetAddCssClass(void *w, const char *c) { gtk_widget_add_css_class((GtkWidget*)w, c); }
void gtk4WidgetSetName(void *w, const char *n) { gtk_widget_set_name((GtkWidget*)w, n); }
void gtk4WidgetSetTooltip(void *w, const char *t) { gtk_widget_set_tooltip_text((GtkWidget*)w, t); }

void gtk4SetDarkTheme(int dark) {
	GtkSettings *settings = gtk_settings_get_default();
	g_object_set(settings, "gtk-application-prefer-dark-theme", dark ? TRUE : FALSE, NULL);
}

void gtk4SetThemeName(const char *name) {
	GtkSettings *settings = gtk_settings_get_default();
	g_object_set(settings, "gtk-theme-name", name, NULL);
}

void gtk4SetIconThemeName(const char *name) {
	GtkSettings *settings = gtk_settings_get_default();
	g_object_set(settings, "gtk-icon-theme-name", name, NULL);
}

void *gtk4ComboBoxTextNew(void) { return gtk_combo_box_text_new(); }
void *gtk4ComboBoxTextNewWithEntry(void) { return gtk_combo_box_text_new_with_entry(); }
void gtk4ComboBoxTextAppend(void *cb, const char *id, const char *text) { gtk_combo_box_text_append((GtkComboBoxText*)cb, id, text); }
void gtk4ComboBoxTextSetActive(void *cb, int idx) { gtk_combo_box_set_active((GtkComboBox*)cb, idx); }
int gtk4ComboBoxTextGetActive(void *cb) { return gtk_combo_box_get_active((GtkComboBox*)cb); }
const char *gtk4ComboBoxTextGetActiveId(void *cb) { return gtk_combo_box_get_active_id((GtkComboBox*)cb); }
char *gtk4ComboBoxTextGetActiveText(void *cb) { return gtk_combo_box_text_get_active_text((GtkComboBoxText*)cb); }
void gtk4ComboBoxTextRemoveAll(void *cb) { gtk_combo_box_text_remove_all((GtkComboBoxText*)cb); }

void *gtk4CheckButtonNew(void) { return gtk_check_button_new(); }
void *gtk4CheckButtonNewWithLabel(const char *label) { return gtk_check_button_new_with_label(label); }
void gtk4CheckButtonSetActive(void *cb, int v) { gtk_check_button_set_active((GtkCheckButton*)cb, v); }
int gtk4CheckButtonGetActive(void *cb) { return gtk_check_button_get_active((GtkCheckButton*)cb); }

void *gtk4SpinButtonNew(double min, double max, double step) {
	GtkAdjustment *adj = gtk_adjustment_new(0, min, max, step, step*10, 0);
	return gtk_spin_button_new(adj, step, 0);
}
void gtk4SpinButtonSetValue(void *sb, double v) { gtk_spin_button_set_value((GtkSpinButton*)sb, v); }
double gtk4SpinButtonGetValue(void *sb) { return gtk_spin_button_get_value((GtkSpinButton*)sb); }

void gtk4WindowSetModal(void *win, int modal) { gtk_window_set_modal((GtkWindow*)win, modal); }
void gtk4WindowSetTransientFor(void *win, void *parent) { gtk_window_set_transient_for((GtkWindow*)win, (GtkWindow*)parent); }

extern void goBridgeIdle(uintptr_t data);

static int _bridge_idle(gpointer data) {
	goBridgeIdle((uintptr_t)data);
	return 0;
}

unsigned int gtk4IdleAdd(uintptr_t data) {
	return g_idle_add(_bridge_idle, (gpointer)(uintptr_t)data);
}

void *gtk4StackSwitcherNew(void) { return gtk_stack_switcher_new(); }
void gtk4StackSwitcherSetStack(void *sw, void *stack) { gtk_stack_switcher_set_stack((GtkStackSwitcher*)sw, (GtkStack*)stack); }
void *gtk4WindowGetTransientFor(void *win) { return gtk_window_get_transient_for((GtkWindow*)win); }
