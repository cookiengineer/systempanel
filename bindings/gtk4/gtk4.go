package gtk4

/*
#cgo pkg-config: gtk4
#cgo CFLAGS: -Wno-deprecated-declarations

#include <stdint.h>
#include <stdlib.h>

extern void connectGSignal(void *instance, const char *signal, int handlerType, uintptr_t handle);
extern void *gtk4Init(void);
extern void *gtk4AppNew(const char *id);
extern int gtk4AppRun(void *app, int argc, char **argv);
extern void gtk4AppQuit(void *app);
extern void gtk4AppAddWindow(void *app, void *win);
extern void *gtk4WindowNew(void);
extern void *gtk4AppWindowNew(void *app);
extern void gtk4WindowSetTitle(void *win, const char *title);
extern void gtk4WindowSetDefaultSize(void *win, int w, int h);
extern void gtk4WindowSetChild(void *win, void *child);
extern void gtk4WindowSetTitlebar(void *win, void *bar);
extern void gtk4WindowPresent(void *win);
extern void gtk4WindowClose(void *win);
extern void gtk4WindowDestroy(void *win);
extern void *gtk4HeaderBarNew(void);
extern void gtk4HeaderBarSetTitleWidget(void *bar, void *w);
extern void gtk4HeaderBarPackStart(void *bar, void *w);
extern void gtk4HeaderBarPackEnd(void *bar, void *w);
extern void gtk4HeaderBarSetShowTitleButtons(void *bar, int show);
extern void *gtk4BoxNew(int orientation, int spacing);
extern void gtk4BoxAppend(void *box, void *child);
extern void gtk4BoxPrepend(void *box, void *child);
extern void gtk4BoxRemove(void *box, void *child);
extern void gtk4BoxSetSpacing(void *box, int spacing);
extern void gtk4BoxSetHomogeneous(void *box, int v);
extern void *gtk4ButtonNew(void);
extern void *gtk4ButtonNewWithLabel(const char *label);
extern void gtk4ButtonSetLabel(void *btn, const char *label);
extern void gtk4ButtonSetIconName(void *btn, const char *icon);
extern void gtk4ButtonSetChild(void *btn, void *child);
extern void *gtk4LabelNew(const char *text);
extern void gtk4LabelSetText(void *lbl, const char *text);
extern void gtk4LabelSetMarkup(void *lbl, const char *markup);
extern void gtk4LabelSetWrap(void *lbl, int v);
extern void gtk4LabelSetXAlign(void *lbl, float a);
extern const char *gtk4LabelGetText(void *lbl);
extern void *gtk4ImageNewFromIconName(const char *name);
extern void gtk4ImageSetIconName(void *img, const char *name);
extern void gtk4ImageSetPixelSize(void *img, int size);
extern void *gtk4StackNew(void);
extern void *gtk4StackAddTitled(void *s, void *child, const char *name, const char *title);
extern void gtk4StackSetVisibleChildName(void *s, const char *name);
extern const char *gtk4StackGetVisibleChildName(void *s);
extern void *gtk4StackGetChildByName(void *s, const char *name);
extern void gtk4StackSetTransitionType(void *s, int t);
extern void gtk4StackSetTransitionDuration(void *s, unsigned int d);
extern void gtk4StackSetVHomogeneous(void *s, int v);
extern void *gtk4StackGetPage(void *s, void *child);
extern void gtk4StackPageSetTitle(void *page, const char *title);
extern void gtk4StackPageSetIconName(void *page, const char *name);
extern void gtk4StackPageSetVisible(void *page, int v);
extern void *gtk4ListBoxNew(void);
extern void gtk4ListBoxAppend(void *lb, void *row);
extern void gtk4ListBoxRemove(void *lb, void *row);
extern void gtk4ListBoxSelectRow(void *lb, void *row);
extern void *gtk4ListBoxGetSelectedRow(void *lb);
extern void gtk4ListBoxSetSelectionMode(void *lb, int mode);
extern void *gtk4ListBoxRowNew(void);
extern void gtk4ListBoxRowSetChild(void *row, void *child);
extern int gtk4ListBoxRowGetIndex(void *row);
extern void *gtk4ScrolledWindowNew(void);
extern void gtk4ScrolledWindowSetChild(void *sw, void *child);
extern void gtk4ScrolledWindowSetPolicy(void *sw, int h, int v);
extern void *gtk4ScaleNew(int orientation);
extern void *gtk4ScaleNewWithRange(int orientation, double min, double max, double step);
extern double gtk4ScaleGetValue(void *s);
extern void gtk4ScaleSetValue(void *s, double v);
extern void gtk4ScaleSetRange(void *s, double min, double max);
extern void *gtk4SwitchNew(void);
extern void gtk4SwitchSetActive(void *sw, int v);
extern int gtk4SwitchGetActive(void *sw);
extern void *gtk4EntryNew(void);
extern void gtk4EntrySetText(void *e, const char *text);
extern const char *gtk4EntryGetText(void *e);
extern void gtk4EntrySetPlaceholder(void *e, const char *text);
extern void *gtk4CssProviderNew(void);
extern void gtk4CssLoadFromString(void *css, const char *data);
extern void gtk4CssApplyToDisplay(void *css, unsigned int priority);
extern void *gtk4SpinnerNew(void);
extern void gtk4SpinnerStart(void *s);
extern void gtk4SpinnerStop(void *s);
extern void *gtk4LevelBarNew(void);
extern void gtk4LevelBarSetValue(void *lb, double v);
extern void *gtk4ControllerKeyNew(void);
extern void gtk4WidgetAddController(void *w, void *ctrl);
extern void gtk4WidgetShow(void *w);
extern void gtk4WidgetHide(void *w);
extern void gtk4WidgetSetSensitive(void *w, int v);
extern void gtk4WidgetSetSizeRequest(void *w, int width, int height);
extern void gtk4WidgetSetHExpand(void *w, int v);
extern void gtk4WidgetSetVExpand(void *w, int v);
extern void gtk4WidgetSetHAlign(void *w, int a);
extern void gtk4WidgetSetVAlign(void *w, int a);
extern void gtk4WidgetSetMarginStart(void *w, int m);
extern void gtk4WidgetSetMarginEnd(void *w, int m);
extern void gtk4WidgetSetMarginTop(void *w, int m);
extern void gtk4WidgetSetMarginBottom(void *w, int m);
extern void gtk4WidgetAddCssClass(void *w, const char *c);
extern void gtk4WidgetSetName(void *w, const char *n);
extern void gtk4WidgetSetTooltip(void *w, const char *t);

extern void *gtk4ComboBoxTextNew(void);
extern void *gtk4ComboBoxTextNewWithEntry(void);
extern void gtk4ComboBoxTextAppend(void *cb, const char *id, const char *text);
extern void gtk4ComboBoxTextSetActive(void *cb, int idx);
extern int gtk4ComboBoxTextGetActive(void *cb);
extern const char *gtk4ComboBoxTextGetActiveId(void *cb);
extern char *gtk4ComboBoxTextGetActiveText(void *cb);
extern void gtk4ComboBoxTextRemoveAll(void *cb);

extern void *gtk4CheckButtonNew(void);
extern void *gtk4CheckButtonNewWithLabel(const char *label);
extern void gtk4CheckButtonSetActive(void *cb, int v);
extern int gtk4CheckButtonGetActive(void *cb);

extern void *gtk4SpinButtonNew(double min, double max, double step);
extern void gtk4SpinButtonSetValue(void *sb, double v);
extern double gtk4SpinButtonGetValue(void *sb);

extern void gtk4WindowSetModal(void *win, int modal);
extern void gtk4WindowSetTransientFor(void *win, void *parent);
extern unsigned int gtk4IdleAdd(uintptr_t data);

extern void *gtk4StackSwitcherNew(void);
extern void gtk4StackSwitcherSetStack(void *sw, void *stack);
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

func init() { C.gtk4Init() }

const (
	sigVoid = iota
	sigAppActivate
	sigListBoxRow
	sigKeyPressed
	sigCloseRequest
	sigScaleValue
	sigNotify
)

//export goBridgeVoid
func goBridgeVoid(data C.uintptr_t) {
	cgo.Handle(uintptr(data)).Value().(func())()
}

//export goBridgeListBoxRowActivated
func goBridgeListBoxRowActivated(data C.uintptr_t, rowPtr unsafe.Pointer) {
	h := cgo.Handle(uintptr(data))
	row := &ListBoxRow{Widget: Widget{ptr: rowPtr}}
	h.Value().(func(*ListBoxRow))(row)
}

//export goBridgeBool3Uint
func goBridgeBool3Uint(data C.uintptr_t, a, b, c C.uint) C.int {
	if cgo.Handle(uintptr(data)).Value().(func(uint,uint,uint) bool)(uint(a), uint(b), uint(c)) {
		return 1
	}
	return 0
}

//export goBridgeBool
func goBridgeBool(data C.uintptr_t) C.int {
	if cgo.Handle(uintptr(data)).Value().(func() bool)() { return 1 }
	return 0
}

//export goBridgeDouble
func goBridgeDouble(data C.uintptr_t, val C.double) {
	cgo.Handle(uintptr(data)).Value().(func(float64))(float64(val))
}

//export goBridgeIdle
func goBridgeIdle(data C.uintptr_t) {
	cgo.Handle(uintptr(data)).Value().(func())()
}

type Widget struct{ ptr unsafe.Pointer }
func (w Widget) Ptr() unsafe.Pointer { return w.ptr }
func (w *Widget) Show()              { C.gtk4WidgetShow(w.ptr) }
func (w *Widget) Hide()              { C.gtk4WidgetHide(w.ptr) }
func (w *Widget) SetSensitive(v bool) { C.gtk4WidgetSetSensitive(w.ptr, cbool(v)) }
func (w *Widget) SetSizeRequest(wd, ht int) { C.gtk4WidgetSetSizeRequest(w.ptr, C.int(wd), C.int(ht)) }
func (w *Widget) SetHExpand(v bool)  { C.gtk4WidgetSetHExpand(w.ptr, cbool(v)) }
func (w *Widget) SetVExpand(v bool)  { C.gtk4WidgetSetVExpand(w.ptr, cbool(v)) }
func (w *Widget) SetHAlign(a Align)  { C.gtk4WidgetSetHAlign(w.ptr, C.int(a)) }
func (w *Widget) SetVAlign(a Align)  { C.gtk4WidgetSetVAlign(w.ptr, C.int(a)) }
func (w *Widget) SetMarginStart(m int)  { C.gtk4WidgetSetMarginStart(w.ptr, C.int(m)) }
func (w *Widget) SetMarginEnd(m int)    { C.gtk4WidgetSetMarginEnd(w.ptr, C.int(m)) }
func (w *Widget) SetMarginTop(m int)    { C.gtk4WidgetSetMarginTop(w.ptr, C.int(m)) }
func (w *Widget) SetMarginBottom(m int) { C.gtk4WidgetSetMarginBottom(w.ptr, C.int(m)) }
func (w *Widget) AddCSSClass(c string) { cs := C.CString(c); defer C.free(unsafe.Pointer(cs)); C.gtk4WidgetAddCssClass(w.ptr, cs) }
func (w *Widget) SetName(n string) { cn := C.CString(n); defer C.free(unsafe.Pointer(cn)); C.gtk4WidgetSetName(w.ptr, cn) }
func (w *Widget) SetTooltip(t string) { ct := C.CString(t); defer C.free(unsafe.Pointer(ct)); C.gtk4WidgetSetTooltip(w.ptr, ct) }
func (w *Widget) RemoveCSSClass(c string) {}

type Align int
const (AlignFill Align=0; AlignStart Align=1; AlignEnd Align=2; AlignCenter Align=3; AlignBaseline Align=4)
type Orientation int
const (OrientationHorizontal Orientation=0; OrientationVertical Orientation=1)
type SelectionMode int
const (SelectionNone SelectionMode=0; SelectionSingle SelectionMode=1; SelectionBrowse SelectionMode=2; SelectionMultiple SelectionMode=3)
type StackTransitionType int
const (StackTransitionNone StackTransitionType=0; StackTransitionCrossfade StackTransitionType=1; StackTransitionSlideRight StackTransitionType=2; StackTransitionSlideLeft StackTransitionType=3)
type PolicyType int
const (PolicyAlways PolicyType=0; PolicyAutomatic PolicyType=1; PolicyNever PolicyType=2; PolicyExternal PolicyType=3)
type License int
const (
	CSSPriorityFallback = 0; CSSPriorityTheme = 200; CSSPrioritySettings = 400
	CSSPriorityApplication = 600; CSSPriorityUser = 800
)

func cbool(v bool) C.int { if v { return 1 }; return 0 }

func csig(name string) *C.char { return C.CString(name) }

func gsig(instance unsafe.Pointer, signal string, sigType int, fn interface{}) {
	h := cgo.NewHandle(fn)
	cs := C.CString(signal)
	C.connectGSignal(instance, cs, C.int(sigType), C.uintptr_t(uintptr(h)))
	C.free(unsafe.Pointer(cs))
}

func gsigVoid(instance unsafe.Pointer, signal string, fn func())               { gsig(instance, signal, sigVoid, fn) }
func gsigAppActivate(instance unsafe.Pointer, fn func())                        { gsig(instance, "activate", sigAppActivate, fn) }
func gsigListBoxRow(instance unsafe.Pointer, signal string, fn func(*ListBoxRow)) { gsig(instance, signal, sigListBoxRow, fn) }
func gsigBool(instance unsafe.Pointer, signal string, fn func() bool)            { gsig(instance, signal, sigCloseRequest, fn) }
func gsigBool3(instance unsafe.Pointer, signal string, fn func(uint,uint,uint) bool) { gsig(instance, signal, sigKeyPressed, fn) }
func gsigScale(instance unsafe.Pointer, signal string, fn func(float64))          { gsig(instance, signal, sigScaleValue, fn) }

type Application struct{ Widget }
func ApplicationNew(id string) *Application {
	cid := C.CString(id); defer C.free(unsafe.Pointer(cid))
	return &Application{Widget{ptr: C.gtk4AppNew(cid)}}
}
func (a *Application) OnActivate(fn func()) { gsigAppActivate(a.ptr, fn) }
func (a *Application) Run(args []string) int {
	argc := C.int(len(args)); argv := make([]*C.char, len(args))
	for i, a := range args { argv[i] = C.CString(a); defer C.free(unsafe.Pointer(argv[i])) }
	return int(C.gtk4AppRun(a.ptr, argc, &argv[0]))
}
func (a *Application) Quit() { C.gtk4AppQuit(a.ptr) }
func (a *Application) AddWindow(w *Window) { C.gtk4AppAddWindow(a.ptr, w.ptr) }

type Window struct{ Widget }
func WindowNew() *Window { return &Window{Widget{ptr: C.gtk4WindowNew()}} }
func ApplicationWindowNew(app *Application) *Window { return &Window{Widget{ptr: C.gtk4AppWindowNew(app.ptr)}} }
func (w *Window) SetTitle(t string) { ct:=C.CString(t); defer C.free(unsafe.Pointer(ct)); C.gtk4WindowSetTitle(w.ptr, ct) }
func (w *Window) SetDefaultSize(wd, ht int) { C.gtk4WindowSetDefaultSize(w.ptr, C.int(wd), C.int(ht)) }
func (w *Window) SetChild(child *Widget) { C.gtk4WindowSetChild(w.ptr, child.ptr) }
func (w *Window) SetTitlebar(bar *HeaderBar) { C.gtk4WindowSetTitlebar(w.ptr, bar.ptr) }
func (w *Window) Present() { C.gtk4WindowPresent(w.ptr) }
func (w *Window) Close() { C.gtk4WindowClose(w.ptr) }
func (w *Window) Destroy() { C.gtk4WindowDestroy(w.ptr) }
func (w *Window) OnCloseRequest(fn func() bool) { gsigBool(w.ptr, "close-request", fn) }
func (w *Window) SetModal(v bool) { C.gtk4WindowSetModal(w.ptr, cbool(v)) }
func (w *Window) SetTransientFor(parent *Window) { C.gtk4WindowSetTransientFor(w.ptr, parent.ptr) }

type HeaderBar struct{ Widget }
func HeaderBarNew() *HeaderBar { return &HeaderBar{Widget{ptr: C.gtk4HeaderBarNew()}} }
func (h *HeaderBar) SetTitleWidget(w *Widget) { C.gtk4HeaderBarSetTitleWidget(h.ptr, w.ptr) }
func (h *HeaderBar) PackStart(child *Widget) { C.gtk4HeaderBarPackStart(h.ptr, child.ptr) }
func (h *HeaderBar) PackEnd(child *Widget) { C.gtk4HeaderBarPackEnd(h.ptr, child.ptr) }
func (h *HeaderBar) SetShowTitleButtons(v bool) { C.gtk4HeaderBarSetShowTitleButtons(h.ptr, cbool(v)) }

type Box struct{ Widget }
func BoxNew(o Orientation, s int) *Box { return &Box{Widget{ptr: C.gtk4BoxNew(C.int(o), C.int(s))}} }
func (b *Box) Append(child *Widget) { C.gtk4BoxAppend(b.ptr, child.ptr) }
func (b *Box) Prepend(child *Widget) { C.gtk4BoxPrepend(b.ptr, child.ptr) }
func (b *Box) Remove(child *Widget) { C.gtk4BoxRemove(b.ptr, child.ptr) }
func (b *Box) SetSpacing(s int) { C.gtk4BoxSetSpacing(b.ptr, C.int(s)) }
func (b *Box) SetHomogeneous(v bool) { C.gtk4BoxSetHomogeneous(b.ptr, cbool(v)) }

type Button struct{ Widget }
func ButtonNew() *Button { return &Button{Widget{ptr: C.gtk4ButtonNew()}} }
func ButtonNewWithLabel(l string) *Button { cl:=C.CString(l); defer C.free(unsafe.Pointer(cl)); return &Button{Widget{ptr: C.gtk4ButtonNewWithLabel(cl)}} }
func (b *Button) SetLabel(l string) { cl:=C.CString(l); defer C.free(unsafe.Pointer(cl)); C.gtk4ButtonSetLabel(b.ptr, cl) }
func (b *Button) SetIconName(n string) { cn:=C.CString(n); defer C.free(unsafe.Pointer(cn)); C.gtk4ButtonSetIconName(b.ptr, cn) }
func (b *Button) SetChild(child *Widget) { C.gtk4ButtonSetChild(b.ptr, child.ptr) }
func (b *Button) OnClicked(fn func()) { gsigVoid(b.ptr, "clicked", fn) }

type Label struct{ Widget }
func LabelNew(t string) *Label { ct:=C.CString(t); defer C.free(unsafe.Pointer(ct)); return &Label{Widget{ptr: C.gtk4LabelNew(ct)}} }
func (l *Label) SetText(t string) { ct:=C.CString(t); defer C.free(unsafe.Pointer(ct)); C.gtk4LabelSetText(l.ptr, ct) }
func (l *Label) SetMarkup(m string) { cm:=C.CString(m); defer C.free(unsafe.Pointer(cm)); C.gtk4LabelSetMarkup(l.ptr, cm) }
func (l *Label) SetWrap(v bool) { C.gtk4LabelSetWrap(l.ptr, cbool(v)) }
func (l *Label) SetXAlign(a float64) { C.gtk4LabelSetXAlign(l.ptr, C.float(a)) }
func (l *Label) Text() string { return C.GoString(C.gtk4LabelGetText(l.ptr)) }

type Image struct{ Widget }
func ImageNewFromIconName(n string) *Image { cn:=C.CString(n); defer C.free(unsafe.Pointer(cn)); return &Image{Widget{ptr: C.gtk4ImageNewFromIconName(cn)}} }
func (i *Image) SetIconName(n string) { cn:=C.CString(n); defer C.free(unsafe.Pointer(cn)); C.gtk4ImageSetIconName(i.ptr, cn) }
func (i *Image) SetPixelSize(s int) { C.gtk4ImageSetPixelSize(i.ptr, C.int(s)) }

type Stack struct{ Widget }
func StackNew() *Stack { return &Stack{Widget{ptr: C.gtk4StackNew()}} }
func (s *Stack) AddTitled(child *Widget, name, title string) *StackPage {
	cn:=C.CString(name); defer C.free(unsafe.Pointer(cn))
	ct:=C.CString(title); defer C.free(unsafe.Pointer(ct))
	return &StackPage{ptr: C.gtk4StackAddTitled(s.ptr, child.ptr, cn, ct)}
}
func (s *Stack) SetVisibleChildName(name string) { cn:=C.CString(name); defer C.free(unsafe.Pointer(cn)); C.gtk4StackSetVisibleChildName(s.ptr, cn) }
func (s *Stack) GetVisibleChildName() string { return C.GoString(C.gtk4StackGetVisibleChildName(s.ptr)) }
func (s *Stack) GetChildByName(name string) *Widget {
	cn:=C.CString(name); defer C.free(unsafe.Pointer(cn)); p := C.gtk4StackGetChildByName(s.ptr, cn)
	if p == nil { return nil }; return &Widget{ptr: p}
}
func (s *Stack) SetTransitionType(t StackTransitionType) { C.gtk4StackSetTransitionType(s.ptr, C.int(t)) }
func (s *Stack) SetTransitionDuration(d uint) { C.gtk4StackSetTransitionDuration(s.ptr, C.uint(d)) }
func (s *Stack) SetVHomogeneous(v bool) { C.gtk4StackSetVHomogeneous(s.ptr, cbool(v)) }
func (s *Stack) SetHHomogeneous(v bool) {}
func (s *Stack) GetPage(child *Widget) *StackPage {
	p := C.gtk4StackGetPage(s.ptr, child.ptr); if p == nil { return nil }; return &StackPage{ptr: p}
}

type StackPage struct{ ptr unsafe.Pointer }
func (p *StackPage) SetTitle(t string) { ct:=C.CString(t); defer C.free(unsafe.Pointer(ct)); C.gtk4StackPageSetTitle(p.ptr, ct) }
func (p *StackPage) SetIconName(n string) { cn:=C.CString(n); defer C.free(unsafe.Pointer(cn)); C.gtk4StackPageSetIconName(p.ptr, cn) }
func (p *StackPage) SetVisible(v bool) { C.gtk4StackPageSetVisible(p.ptr, cbool(v)) }

type ListBox struct{ Widget }
func ListBoxNew() *ListBox { return &ListBox{Widget{ptr: C.gtk4ListBoxNew()}} }
func (lb *ListBox) Append(row *ListBoxRow) { C.gtk4ListBoxAppend(lb.ptr, row.ptr) }
func (lb *ListBox) Remove(row *ListBoxRow) { C.gtk4ListBoxRemove(lb.ptr, row.ptr) }
func (lb *ListBox) SelectRow(row *ListBoxRow) { C.gtk4ListBoxSelectRow(lb.ptr, row.ptr) }
func (lb *ListBox) GetSelectedRow() *ListBoxRow {
	r := C.gtk4ListBoxGetSelectedRow(lb.ptr); if r == nil { return nil }; return &ListBoxRow{Widget{ptr: r}}
}
func (lb *ListBox) SetSelectionMode(m SelectionMode) { C.gtk4ListBoxSetSelectionMode(lb.ptr, C.int(m)) }
func (lb *ListBox) OnRowActivated(fn func(*ListBoxRow)) { gsigListBoxRow(lb.ptr, "row-selected", fn) }

type ListBoxRow struct{ Widget }
func ListBoxRowNew() *ListBoxRow { return &ListBoxRow{Widget{ptr: C.gtk4ListBoxRowNew()}} }
func (r *ListBoxRow) SetChild(child *Widget) { C.gtk4ListBoxRowSetChild(r.ptr, child.ptr) }
func (r *ListBoxRow) GetIndex() int { return int(C.gtk4ListBoxRowGetIndex(r.ptr)) }
func (r *ListBoxRow) SetActivatable(v bool) {}
func (r *ListBoxRow) SetSelectable(v bool) {}

type ScrolledWindow struct{ Widget }
func ScrolledWindowNew() *ScrolledWindow { return &ScrolledWindow{Widget{ptr: C.gtk4ScrolledWindowNew()}} }
func (sw *ScrolledWindow) SetChild(child *Widget) { C.gtk4ScrolledWindowSetChild(sw.ptr, child.ptr) }
func (sw *ScrolledWindow) SetPolicy(h, v PolicyType) { C.gtk4ScrolledWindowSetPolicy(sw.ptr, C.int(h), C.int(v)) }
func (sw *ScrolledWindow) SetPropagateNaturalHeight(v bool) {}
func (sw *ScrolledWindow) SetPropagateNaturalWidth(v bool) {}

type Scale struct{ Widget }
func ScaleNew(o Orientation) *Scale { return &Scale{Widget{ptr: C.gtk4ScaleNew(C.int(o))}} }
func ScaleNewWithRange(o Orientation, min, max, step float64) *Scale {
	return &Scale{Widget{ptr: C.gtk4ScaleNewWithRange(C.int(o), C.double(min), C.double(max), C.double(step))}}
}
func (s *Scale) SetValue(v float64) { C.gtk4ScaleSetValue(s.ptr, C.double(v)) }
func (s *Scale) GetValue() float64 { return float64(C.gtk4ScaleGetValue(s.ptr)) }
func (s *Scale) SetRange(min, max float64) { C.gtk4ScaleSetRange(s.ptr, C.double(min), C.double(max)) }
func (s *Scale) OnValueChanged(fn func(float64)) { gsigScale(s.ptr, "value-changed", fn) }

type Switch struct{ Widget }
func SwitchNew() *Switch { return &Switch{Widget{ptr: C.gtk4SwitchNew()}} }
func (sw *Switch) SetActive(v bool) { C.gtk4SwitchSetActive(sw.ptr, cbool(v)) }
func (sw *Switch) GetActive() bool { return C.gtk4SwitchGetActive(sw.ptr) != 0 }
func (sw *Switch) SetState(v bool) { sw.SetActive(v) }
func (sw *Switch) GetState() bool { return sw.GetActive() }
func (sw *Switch) OnActivate(fn func()) { gsigVoid(sw.ptr, "activate", fn) }

type Entry struct{ Widget }
func EntryNew() *Entry { return &Entry{Widget{ptr: C.gtk4EntryNew()}} }
func (e *Entry) SetText(t string) { ct:=C.CString(t); defer C.free(unsafe.Pointer(ct)); C.gtk4EntrySetText(e.ptr, ct) }
func (e *Entry) GetText() string { return C.GoString(C.gtk4EntryGetText(e.ptr)) }
func (e *Entry) SetPlaceholder(t string) { ct:=C.CString(t); defer C.free(unsafe.Pointer(ct)); C.gtk4EntrySetPlaceholder(e.ptr, ct) }
func (e *Entry) SetVisibility(v bool) {}
func (e *Entry) OnChanged(fn func()) { gsigVoid(e.ptr, "changed", fn) }

type CSS struct{ ptr unsafe.Pointer }
func CSSNew() *CSS { return &CSS{ptr: C.gtk4CssProviderNew()} }
func (c *CSS) LoadFromString(css string) { cc:=C.CString(css); defer C.free(unsafe.Pointer(cc)); C.gtk4CssLoadFromString(c.ptr, cc) }
func (c *CSS) ApplyToDisplay(priority uint) { C.gtk4CssApplyToDisplay(c.ptr, C.uint(priority)) }

type Spinner struct{ Widget }
func SpinnerNew() *Spinner { return &Spinner{Widget{ptr: C.gtk4SpinnerNew()}} }
func (s *Spinner) Start() { C.gtk4SpinnerStart(s.ptr) }
func (s *Spinner) Stop() { C.gtk4SpinnerStop(s.ptr) }
func (s *Spinner) SetSpinning(v bool) { if v { s.Start() } else { s.Stop() } }

type LevelBar struct{ Widget }
func LevelBarNew() *LevelBar { return &LevelBar{Widget{ptr: C.gtk4LevelBarNew()}} }
func (lb *LevelBar) SetValue(v float64) { C.gtk4LevelBarSetValue(lb.ptr, C.double(v)) }

type EventControllerKey struct{ Widget }
func EventControllerKeyNew() *EventControllerKey { return &EventControllerKey{Widget{ptr: C.gtk4ControllerKeyNew()}} }
func (e *EventControllerKey) OnKeyPressed(fn func(uint,uint,uint) bool) { gsigBool3(e.ptr, "key-pressed", fn) }
func (e *EventControllerKey) AddToWidget(w *Widget) { C.gtk4WidgetAddController(w.ptr, e.ptr) }

type ComboBoxText struct{ Widget }
func ComboBoxTextNew() *ComboBoxText { return &ComboBoxText{Widget{ptr: C.gtk4ComboBoxTextNew()}} }
func ComboBoxTextNewWithEntry() *ComboBoxText { return &ComboBoxText{Widget{ptr: C.gtk4ComboBoxTextNewWithEntry()}} }
func (c *ComboBoxText) Append(id, text string) {
	cid:=C.CString(id); defer C.free(unsafe.Pointer(cid))
	ct:=C.CString(text); defer C.free(unsafe.Pointer(ct))
	C.gtk4ComboBoxTextAppend(c.ptr, cid, ct)
}
func (c *ComboBoxText) SetActive(idx int) { C.gtk4ComboBoxTextSetActive(c.ptr, C.int(idx)) }
func (c *ComboBoxText) GetActive() int { return int(C.gtk4ComboBoxTextGetActive(c.ptr)) }
func (c *ComboBoxText) GetActiveID() string { return C.GoString(C.gtk4ComboBoxTextGetActiveId(c.ptr)) }
func (c *ComboBoxText) GetActiveText() string {
	s := C.gtk4ComboBoxTextGetActiveText(c.ptr)
	if s == nil { return "" }
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}
func (c *ComboBoxText) RemoveAll() { C.gtk4ComboBoxTextRemoveAll(c.ptr) }

type CheckButton struct{ Widget }
func CheckButtonNew() *CheckButton { return &CheckButton{Widget{ptr: C.gtk4CheckButtonNew()}} }
func CheckButtonNewWithLabel(l string) *CheckButton {
	cl:=C.CString(l); defer C.free(unsafe.Pointer(cl))
	return &CheckButton{Widget{ptr: C.gtk4CheckButtonNewWithLabel(cl)}}
}
func (cb *CheckButton) SetActive(v bool) { C.gtk4CheckButtonSetActive(cb.ptr, cbool(v)) }
func (cb *CheckButton) GetActive() bool { return C.gtk4CheckButtonGetActive(cb.ptr) != 0 }

type SpinButton struct{ Widget }
func SpinButtonNew(min, max, step float64) *SpinButton {
	return &SpinButton{Widget{ptr: C.gtk4SpinButtonNew(C.double(min), C.double(max), C.double(step))}}
}
func (s *SpinButton) SetValue(v float64) { C.gtk4SpinButtonSetValue(s.ptr, C.double(v)) }
func (s *SpinButton) GetValue() float64 { return float64(C.gtk4SpinButtonGetValue(s.ptr)) }

type StackSwitcher struct{ Widget }
func StackSwitcherNew() *StackSwitcher { return &StackSwitcher{Widget{ptr: C.gtk4StackSwitcherNew()}} }
func (sw *StackSwitcher) SetStack(s *Stack) { C.gtk4StackSwitcherSetStack(sw.ptr, s.ptr) }

func IdleAdd(fn func()) {
	h := cgo.NewHandle(fn)
	C.gtk4IdleAdd(C.uintptr_t(uintptr(h)))
}
