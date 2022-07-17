package dom

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)
var(
	DEFAULT = Window{size:Size(300,300),pos:Pos(100,100)}
	upgrader = websocket.Upgrader{}
	port string = ":3333"
	rootws string = "/ws"
	contenido string  = "<html><head></head></html>"
	rootserv string = "http://localhost" + port + "/"
	com *websocket.Conn
	conection bool = false
	done bool = false
	methods []method
	Event *Events
	Dom []*Element
	document *Element = &Element{ TagName: "document",}
	childsApp []*Component
	states []*State
	tags string = "img,section,div,h1,h2,h3,h4,h5,h6,input,label,button,body,comp"
)
// types
type State struct {
	name string
	value string
	parent *Element
}
type method struct{
   name string
   function func()
}
type size struct{
	Width int
	Height int
}
type pos struct {
	PosX int
	PosY int
}
type Window struct {
	size size
	pos pos
	icon string
	title string
}
type smsJsons struct{
   Type string `json:"type"`
   Ref string `json:"ref"`
   Body string `json:"body"`
   Value string `json:"value"`
   Field string `json:"field"`
   Event string `json:"event"`
   Name string `json:"name"`
}
type Component struct{
	name string 
	Model func()string
	Action func()
}
type Element struct {
	ref int
	TagName    string
	innerHtml  string
	OuterHtml  string
	Value      string
	ClassName string
	Id      string
	Name string
	ParentNode *Element
	Children   []*Element
	parentComponent *Component
}
type Events struct{
	Type string `json:"type"`
	Value string `json:"value"`
	Ref string `json:"ref"`
}
type Attrs struct{
	Type string
	Value string
}
// window 
func NewWindow()*Window{
	return &Window{Size(300,300),Pos(100,100),"",""}
}
func (win *Window)SetSize(w int , h int )*Window{
	win.size.Width = w
	win.size.Height = h
	return win
}
func (w *Window)SetPosition(p pos)*Window{
	w.pos = p
	return w
}
func (w *Window)SetTitle(t string)*Window{
	w.title = t
	return w
}
func (w *Window)SetIcon(i string)*Window{
	w.icon = i
	return w
}
func (w *Window)PositionCenter() *Window{
	w.pos = Center()
	return w
}
func Size(w int , h int)size{
	return size{Width: w,Height: h}
}
func (p pos) dividir(c int)pos{
	p.PosX = p.PosX/c
	p.PosY = p.PosY/c
	return p
} 
func Pos(x int , y int)pos{
	return pos{PosX: x,PosY: y}
}
func Center()pos{
	return Pos(1300,800).dividir(4)
}
func SizeDefault()size{
	return Size(300,300)
}
// utils 
func Error(err error){
	if err != nil{
		fmt.Println(err)
	}
}
func Log( s ...interface{}){
	fmt.Println( s... )
}
func ToFirstUpperCase(str string)string{
	str = strings.ToTitle(str[:1]) + str[1:]
	return str
}
func Clean(s string)string{
	res := strings.ReplaceAll(strings.ReplaceAll(s , "\t",""),"\n","")
	if strings.Contains(res ,"> <"){
		res = strings.ReplaceAll(res ,"> <","><")
	}
	return res
}
func (p *Element) uploadInners(){
	parent := p
	for{
		if parent.ParentNode != nil{
			parent = parent.ParentNode
			parent.innerHtml = ""
			endCabecera := strings.Index(parent.OuterHtml,">")
			open := parent.OuterHtml[:endCabecera+1]
			close := "</"+ parent.TagName + ">"
			for _, v := range parent.Children{
				parent.innerHtml += v.OuterHtml
			}
			parent.OuterHtml = open + parent.innerHtml + close
		}else{ break }
	}
}
func getNameByMethod(f string,m string)(name string){
	file,err:= os.ReadFile(f)
	if err != nil { fmt.Println(err) }
	sliceString := strings.Split(string(file),"\n")
	for _, line := range sliceString{
		if strings.Contains(line,m){
			name = strings.TrimSpace(strings.ReplaceAll(strings.Split(line,"=")[0],":",""))
			if !stateExists(name){
				return 
			}
		}
	}
	return 
}
// eventos
func onWindowLoad(call func()){
	for{
		if conection{
			call()
			return 
		}
	}
}
func OnWait(){
	for {
		if done == true {
			Log("- Closing app ... bye")
			os.Remove("./src/index.html")
			return 
		}
	}
}
// control window
func isWindows(w *Window) bool {

	cmd := exec.Command("C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe", "--app="+fmt.Sprintf("%s", rootserv),"--window-size="+fmt.Sprint(w.size.Width)+","+fmt.Sprint(w.size.Height),"--window-position="+fmt.Sprint(w.pos.PosX)+","+fmt.Sprint(w.pos.PosY))
	if err := cmd.Start(); err != nil { // Ejecutar comando
		cmd := exec.Command("c:/Program Files (x86)/Google/Chrome/Application/./chrome", "--app="+fmt.Sprintf("%s", rootserv),"--window-size="+fmt.Sprint(w.size.Width)+","+fmt.Sprint(w.size.Height),"--window-position="+fmt.Sprint(w.pos.PosX)+","+fmt.Sprint(w.pos.PosY))
		if err := cmd.Start(); err != nil { // Ejecutar comando

			log.Fatal(err)
			return false
		}
	}

	return true
}
func isLinux(w *Window)bool{

	user :=os.Getenv("LOGNAME")

	if err := exec.Command("mkdir","/home/"+user+"/.profileFirefox").Run(); err != nil { fmt.Println("- Error create profile folder"," >> ",err) }else{ fmt.Println("- Create profile folder exited ! ") }
	if err := exec.Command("firefox","-CreateProfile","app /home/"+user+"/.profileFirefox").Run(); err != nil{ fmt.Println("- Error createprofile in firefox >> ",err) }else{ fmt.Println("- Create profile firefox exited ! ") }
	if err := exec.Command("mkdir","/home/"+user+"/.profileFirefox/chrome").Run(); err != nil { fmt.Println("- Error create chrome folder >> ",err) }else{ fmt.Println("- Create chrome floder exited !") }
	file, err := os.Create("/home/"+user+"/.profileFirefox/chrome/userChrome.css")
	if err != nil { fmt.Println("- Error create css file >> ",err) }else{ fmt.Println("- Create css file exited !") }
	_, err = file.Write([]byte(`
						@namespace url("http://www.mozilla.org/keymaster/gatekeeper/there.is.only.xul");

						#navigator-toolbox {
							width:      0px; 
							height:     0px;
							opacity:    0  ;
							overflow:   hidden;
					}`))
	if err != nil { fmt.Println("- Error write css file >> ",err) }else{ fmt.Println("- Write css file exited !") }

	if err := exec.Command("firefox",fmt.Sprintf("%s", rootserv),"-p","app","--window-size="+fmt.Sprint(w.size.Width)+","+fmt.Sprint(w.size.Height),"--window-position="+fmt.Sprint(w.pos.PosX)+","+fmt.Sprint(w.pos.PosY)).Start(); err != nil{

		fmt.Println("Command fail !!")
		return false
	}
	return true
}
func js() string{
	js:= `
	<script>

	const ws = new WebSocket ("ws://localhost`+ port + rootws +`");

	ws.onopen = (e)=>{

		console.log("conectado");
		ws.send("ok");

	};

	ws.onmessage = (e)=>{
		console.log(e.data);
		let data = JSON.parse( e.data );

		if ( data.type == "bind" ){
			isBind( data );
			return;
		}
		if ( data.type == "eval" ){
			isEval( data );
			return; 
		}
	};

	function isBind( data ){

		window[data.name] = ()=>{
			if ( event.type != "dragstart"){
				event.preventDefault();
			}
			console.log(event.target.value);
			ws.send( JSON.stringify({type:"event", name:data.name ,event:JSON.stringify({type:event.type,ref:event.target.getAttribute('key'),value:event.target.value})}) );
		};
	};

	function isEval( data ){

		let res = eval( data.js );

		if ( typeof res != "string" ){
			res = JSON.stringify( res );
		}
		if ( data.id ){
			res = JSON.stringify( {id : data.id , body: res} );
		}
		if ( res != undefined ){
			ws.send( res );
		}
	};


	window.addEventListener('beforeunload', (e)=>{
		e.preventDefault();
		ws.send("close");
	});

	window.uploadValue = ( data )=>{
		let res = data;
		let ele = document.querySelector('.' + data.ref );

		ele.addEventListener('change',()=>{
				res.value = ele.value.toString();
				ws.send(JSON.stringify({type:"upload", Ref:data.ref , value : res.value , body :JSON.stringify(res)}));
		});
	};
	</script>
	`
	return js
}
func New(content Component , w *Window){


	app := NewElementL(Build(content.Model()))

	if !strings.Contains(app.OuterHtml,"<html>") || !strings.Contains(app.OuterHtml,"<body>"){
		app.OuterHtml = `<!DOCTYPE html><meta name="theme-color" content="#872e4e"><html><body>`+ app.OuterHtml +`</html></body>`
	}
	
	app.OuterHtml = strings.Replace(app.OuterHtml,"<html>","<html><head><title>"+ w.title +"</title><link rel='icon' href='"+ w.icon +"' sizes='16x16 32x32' type='image/png'></head>", 1)
	
	
	app.OuterHtml = strings.Replace( app.OuterHtml,"<body>","<body><style>"+ readCss() +"</style>", 1)
	
	
	// inyect js
	app.OuterHtml = strings.Replace(app.OuterHtml,"<body>" , "<body>"+ js() , 1)
	contenido = app.OuterHtml
	writeConten(contenido)

	// start server and window
	go newServer()
	system := runtime.GOOS
	fmt.Println("- We are in", system )
	switch system {
		case "windows":
			isWindows(w)
		case "linux":
			isLinux(w)
		case "mac":
	}
	go onWindowLoad(func(){
		// ejecutar action de todos los componentes
		content.Action()
		for _, child := range childsApp{
			child.Action()
		}
		//uploadValues(replaceVar(simplyConten))
	})
}
// building html con componentes 
func comprovateClassName( ele string){

	if  !strings.Contains(strings.Split(ele,">")[0] , "class"){
		panic("Syntax error : El componente padre debe llevar un className igual que nombre del componente como identificacion interna")
	}
}
func Build( ele string )string{

	go comprovateClassName(ele)

	for _, child := range childsApp{
		if strings.Contains(ele,"</"+child.name+">"){
			ele = strings.Replace(ele,"</"+child.name+">",child.Model(),1)
		}
	}
	//fmt.Println(ele.Model)
	return ele
}
func readCss()string{
	css,_ := os.ReadFile("./src/style.css")
	return string(css)
}
func writeConten(html string){
	file,_ := os.Create("./src/index.html")
	file.WriteString(html)
	file.Close()
}
func GetFile(path string)string{
	conten,err := (os.ReadFile(path))
	if err != nil{fmt.Println(err)}
	return Clean(string(conten))
}
// comunication
func send( sms string ){
	//tipo := "undefined"
	//if strings.Contains(sms,"bind"){ tipo = "is funtion" }else
	//if strings.Contains(sms,"addEventListener"){ tipo = "is event" }else
	//if strings.Contains(sms,"eval"){ tipo = "is evaluation" }
	for {
		if conection{
			_ = com.WriteMessage(1,[]byte(sms))
			//Log("message sent : " , sms  ,"   ", tipo)
			return
		}
	}
}
func upload(s string ){
}
// functions server 
func reciver(w http.ResponseWriter, r *http.Request) {
	com, _ = upgrader.Upgrade(w, r, nil)
	defer com.Close()

	for {
		_, tempSMS, _ := com.ReadMessage()
		// Receive message
		evalOptions(string(tempSMS))
	}
}
func serv(w http.ResponseWriter, r *http.Request){
	w.Write([]byte(contenido))
	http.FileServer(http.Dir("./src/"))
}
func newServer() {
	
	http.HandleFunc(rootws, reciver)
	http.Handle("/" , http.FileServer(http.Dir("./src/")))
	log.Fatal(http.ListenAndServe(port, nil))
}
// add methods 
func AddMethod( name string , f func()){
	methods = append(methods, method{ name , f })
}
func evalMethods( sms string )bool{
	var Json smsJsons
	json.Unmarshal([]byte(sms), &Json)
	json.Unmarshal([]byte(Json.Event),&Event)
	for _,v := range methods{
		if v.name == Json.Name {
			v.function()
			return true
		}
	}
	return false
}
// evaluation options reciver sw
func evalOptions(sms string){

	if sms == "ok"{
		ok()
		return
	}
	if sms == "close"{
		close( sms )
		return
	}
	if strings.Contains( sms , "upload"){
		upload( sms )
		return
	}
	evalMethods(sms)
}
func ok(){
	conection = true
	Log("- Conection is exited !")
}
func close( sms string ){
	done = true
}
// bind js
func Bind(name string , f func()){
	eval(`{"type":"bind","name":"`+ name +`"}`)
	AddMethod( name , f )
	return 
}
func eval( s string ){
	for {
		if conection{
			_ = com.WriteMessage(1,[]byte(s))
			//Log("message sent : " , s )
			return
		}
	}
}
// animation 
func getBoince()string{
	return `@keyframes bounceIn{0%{opacity: 0;transform: scale(0.3) translate3d(0,0,0);}50%{opacity: 0.9;transform: scale(1.1);}80%{opacity: 1;transform: scale(0.89);}100%{opacity: 1;transform: scale(1) translate3d(0,0,0);}}`
}
// instanciate
func NewElement(t string) (e *Element) {
	
		var element Element = Element{
			TagName:   t,
			OuterHtml: "<"+ t +" key='"+ fmt.Sprint(len(Dom))+"'></"+ t +">",
			ref:len(Dom),
		}
		Dom = append(Dom,&element)
		e = &element

	if !strings.Contains(tags,t){
		log.Panic("\033[0;49;33m","Syntax error , tagName no permitted change to  >>","\033[5;49;31m",tags,"\033[;0m")
	}
	return
}
func NewElementL(html string)(*Element){

	var id string
	var className string
	var name string
	var tagName string
	var value string
	var innerHtml string
	var ref int

	ref = len(Dom)

	html = pushKey(html, ref )	
	html = Clean(html)

	i := strings.Index(html, "<")
	i2 := strings.Index(html, ">")
	iTag := strings.Index(html, " ")
	if i > iTag || iTag > i2  {
		iTag = i2
	}
	open := html[i:i2+1]
	tagName = html[i+1:iTag]
	attrs := getAttribute(html[iTag:i2])
	//attrs := strings.Split(html[iTag:i2]," ")
	// obtener innerHtml descartando tags del mismo tipo dentro del elemento
	innerHtml = html[i2+1:]
	tagEnd := "</"+ tagName + ">"
	tagInit := "<"+ tagName 
	end := strings.Index(innerHtml,tagEnd)
	init := strings.Index(innerHtml,tagInit)
	iClose := 0
	iTemp := 0

	for strings.Contains(innerHtml,tagEnd) {
			iTemp = strings.Index(innerHtml,tagEnd)+len(tagEnd)
			iClose += iTemp
			innerHtml = innerHtml[iTemp:]
			if init > end{
				break
			}
	}
	//fmt.Println("tagname:",tagName)
	//fmt.Println("attrs:",attrs)
	//fmt.Println(len(attrs))
	//fmt.Println("inner:",innerHtml)
	innerInit := i2+1
	innerEnd := iClose+i2+1-len(tagEnd)
	if innerInit > innerEnd{
		log.Panic("\033[0;49;33m","Syntax error in >> ","\033[5;49;31m",Dom[len(Dom)-1],"\033[;0m")
	}
	innerHtml = html[innerInit:innerEnd]

	// fin de obtencion
	for _,v := range attrs {
		if strings.Contains(v , "id"){
			id = strings.ReplaceAll(v," ","")[4:len(v)-1]
		}
		if strings.Contains(v , "class"){
			className = strings.ReplaceAll(v," ","")[7:len(v)-1]
		}
		if strings.Contains(v , "name"){
			name = strings.ReplaceAll(v," ","")[6:len(v)-1]
		}
		if strings.Contains(v,"value"){
			value = strings.ReplaceAll(v," ","")[7:len(v)-1]
		}
	}

	ele := NewElement(tagName)

	ele.Id = id
	ele.ClassName = className
	ele.Name = name
	ele.Value = value
	ele.innerHtml = innerHtml
	ele.OuterHtml = open + innerHtml + tagEnd
	ele.ref = ref

	for strings.Contains(innerHtml,"<") && strings.Contains(innerHtml,"</"){
		child := NewElementL(innerHtml)
		existe := false
		for _,v := range Dom{
			if v == child{
				existe = true
				break
			}
		}
		if !existe{
			Dom = append(Dom, child)
		}
		ele.Children = append(ele.Children ,child)
		child.ParentNode = ele
		innerHtml = strings.Replace(innerHtml,child.OuterHtml,"",1)
	}
	return ele
}
func pushKey(s string ,  ref int )(res string){

	if !strings.Contains(s ,"key"){
		slice := strings.Split(s ,">")
		for index , item := range slice{
			if strings.Contains(item , "<") && !strings.Contains(item ,"</"){
				slice[index] += fmt.Sprint(" key='",ref,"'") 
				ref++
			}
			slice[index] += ">"
		}
		res =  strings.Join(slice,"")

	}else{ res = s }
	return 
}
func getAttribute( s string)[]string{

	var out bool = true
	s = strings.TrimSpace(s)
	chars := strings.Split(s,"")

	for index , char := range chars {
		if char == "'" { out = !out }
		if char == " " && out {
			chars[index] = ","
		}
	}
	s = strings.Join(chars, "")

	return strings.Split(s,",")
}
// selectors
func Selector(q string)(ele *Element){

	for _,v := range Dom{

		if strings.Contains(q,"#"){
			if v.Id == q[1:]{
				ele = v 
				break
			}
		}else
		if strings.Contains(q ,"."){
			if strings.Contains(v.ClassName,q[1:]){
				ele = v
				break
			}
		}else
		if v.TagName == q {
				ele = v
				break
		}else 
		if strings.Contains(v.OuterHtml , q){
			ele = v
			break
		}
	}
	if ele == nil { ele = &Element{innerHtml: "<nil>",OuterHtml: "<nil>"}}
	return
}
func SelectorId(q string)(ele *Element){

	for _,v := range Dom{
		if v.Id == q { 
			ele = v 
			break
		}
	}
	return
}
func SelectorAll(q string)(ele []*Element){
	for _,v := range Dom{

		if strings.Contains(q,"#"){
			if v.Id == q[1:]{
				ele = append(ele,v) 
			}
		}else
		if strings.Contains(q ,"."){
			if strings.Contains(v.ClassName ,q[1:]){
				ele = append(ele,v) 
			}
		}else
		if v.TagName == q {
			ele = append(ele,v) 
		}else 
		if strings.Contains(v.OuterHtml , q){
			ele = append(ele,v) 
		}
	}
	if ele == nil { ele =append(ele,&Element{innerHtml: "<nil>",OuterHtml: "<nil>"})}
	return 
}
// methods
func (e *Element) GetRef()string{
	return fmt.Sprint(e.ref)
}
func (e *Element) SetInnerHTML(html string) {
	if e.TagName != "input"{
		clouse := "</" + e.TagName + ">"
		e.innerHtml = html
		e.OuterHtml = strings.Replace(e.OuterHtml , clouse , html + clouse, 1)
		js := "document.querySelector(`[key='"+ fmt.Sprint(e.ref) +"']`).innerHTML =`" + html + "`"
		eval( `{"type":"eval","js":"`+ js +`"}` )	}
}
func (e *Element) uploadInnerHTML(html string) {
	if e.TagName != "input"{
		js := "document.querySelector(`[key='"+ fmt.Sprint(e.ref) +"']`).innerHTML =`" + html + "`"
		eval( `{"type":"eval","js":"`+ js +`"}` )	}
}
func (e *Element) GetInnerHTML() string {
	if e.TagName != "input"{
		return e.innerHtml
	}
	return ""
}
func (e *Element) Append(ele *Element) {
	close := "</" + e.TagName + ">"
	open := e.OuterHtml[:strings.Index(e.OuterHtml,">")+1]
	e.Children = append(e.Children, ele)
	e.innerHtml += ele.OuterHtml
	e.OuterHtml = open + e.innerHtml + close
	if ele.ParentNode != nil {
		ele.Remove()
	}
	ele.ParentNode = e
	js := "document.querySelector(`[key='"+e.GetRef()+"']`).appendChild(document.querySelector(`[key='"+ele.GetRef()+"']`))"
	eval( `{"type":"eval","js":"`+ js +`"}` )
}
func (e *Element) Remove(){
	p := e.ParentNode
	p.innerHtml = strings.Replace(p.innerHtml,e.OuterHtml,"",1)
	p.OuterHtml = strings.Replace(p.OuterHtml,e.OuterHtml,"",1)

	for i ,v := range  p.Children{
		if v == e{
			p.Children = append(p.Children[:i],p.Children[i+1:]...)
		}
	}
	p.uploadInners()
}
func (e *Element) SetAttribute( t string, v string ){

	if t == "name"{
		e.SetName(v)
	}
	if t == "className"{
		e.SetClassName(v)
	}
	if t == "id"{
		e.SetId(v)
	}
	if t == "value"{
		e.Value = v
	}
	js := "document.querySelector(`[key='"+ e.GetRef() +"']`).setAttribute('"+ t +"','"+ v +"')"
	eval(`{"type":"eval","js":"`+ js +`"}`)
}
func (e *Element) SetId(v string) {
	e.Id = v 
	e.OuterHtml = strings.Replace(e.OuterHtml,e.TagName , e.TagName + " id='" + v + "'",1)
	js := "document.querySelector(`[key='"+ e.GetRef() +"']`).id = '"+ v +"'"
	eval(`{"type":"eval","js":"`+ js +`"}`)
}
func (e *Element) SetClassName( v string ) {
	e.ClassName = v
	e.OuterHtml = strings.Replace(e.OuterHtml,e.TagName , e.TagName + " class='" + v + "'" ,1)
	js := "document.querySelector(`[key='"+ e.GetRef() +"']`).className = '"+ v +"'"
	eval(`{"type":"eval","js":"`+ js +`"}`)
}
func (e *Element) SetName( v string ) {
	e.Name = v
	e.OuterHtml = strings.Replace(e.OuterHtml,e.TagName , e.TagName + " name='" + v + "'" ,1)
	js := "document.querySelector(`[key='"+ e.GetRef() +"']`).name = '"+ v +"'"
	eval(`{"type":"eval","js":"`+ js +`"}`)
}
func (e *Element) SetValue(v string){
	e.Value = v
	js := "document.querySelector(`[key='"+ e.GetRef() +"']`).value = '"+ v +"'"
	eval(`{"type":"eval","js":"`+ js +`"}`)
}
func (e *Element) AddEventListener(t string, f func()){
	n := "method_"+fmt.Sprint(len(methods))
	Bind(n,f)
	js := "document.querySelector(`[key='"+ e.GetRef()+"']`).addEventListener('"+ t +"',"+ n +")"
	eval (`{"type":"eval","js":"`+ js +`"}`)
}
func (ev *Events) GetTarget()(ele *Element){

	for i, e := range Dom{
		res,_ := strconv.Atoi(ev.Ref)
		if e.ref == res {
			Dom[i].SetValue(ev.Value)
			ele = e
		}
	}
	if ele == nil { ele = &Element{}}
	return
}
// components
func NewComponent(action func(),model func()string)Component{
	_,caller,_,_ := runtime.Caller(1)
	caller = strings.TrimSpace(caller[strings.LastIndex(caller,"/")+1:len(caller)-3])
	return Component{Action:action,Model:model ,name:caller}
}
func AddChilds(childs ...*Component){
	childsApp = append(childsApp,childs...)
}
func (e *Component) AddChilds(childs ...*Component){
	AddChilds(childs...)
}
func (e *Component) SetName(n string){
	e.name = n
}
// states
func NewState( v string)(res *State){
	_,caller,_,_ := runtime.Caller(1) 
	n := getNameByMethod(caller,"dom.NewState")
	caller = strings.ToLower(caller[strings.LastIndex(caller,"/")+1:len(caller)-3])
	for _,ele := range Dom{
		if ele.ClassName == caller{
			states = append(states,&State{name: n,value: v,parent:ele})
			inner := ele.innerHtml
			for _,item := range states{
				if item.parent == ele{
					inner = strings.ReplaceAll(inner,"$"+item.name,item.value)
				}
			}
			ele.uploadInnerHTML(inner)
			res = states[len(states)-1]
		}
	}
	return 
}
func stateExists(n string) bool{
	for _, v := range states{
		if v.name == n {
			return true
		}
	}
	return false
}
func ( s *State ) Get()string{
	return s.value
}
func ( s *State ) Set(v string){
	inner := s.parent.innerHtml
	s.value = v
	for _,item := range states{
				inner = strings.ReplaceAll(inner,"$"+item.name,item.value)
	}
	s.parent.uploadInnerHTML(inner)
}
func Delay(t time.Duration , callback func()){
	go func(){
		time.Sleep(t)
		callback()
	}()
}








