(window.webpackJsonp=window.webpackJsonp||[]).push([[0],{129:function(e,t,a){e.exports=a(258)},253:function(e,t,a){},255:function(e,t,a){},256:function(e,t,a){},258:function(e,t,a){"use strict";a.r(t);var n=a(0),l=a.n(n),r=a(11),c=a.n(r),i=a(301),s=a(16),o=a(17),u=a(19),m=a(18),d=a(20),p=a(292),h=a(300),f=a(109),E=function(e){function t(e){var a;Object(s.a)(this,t),a=Object(u.a)(this,Object(m.a)(t).call(this,e));for(var n=new Date,l=[],r=0;r<a.props.chartData.length;r++)n.setDate(n.getDate()-r),l.unshift(n.getUTCMonth()+1+"/"+n.getDate()),n=new Date;return a.state={chartData:{labels:l,datasets:[{label:"Outages",backgroundColor:"rgba(255, 99, 132, 0.2)",borderColor:"rgba(255,99,132,1)",borderWidth:2,hoverBackgroundColor:"rgba(255,99,132,0.4)",hoverBorderColor:"rgba(255,99,132,1)",data:a.props.chartData},{label:"404",backgroundColor:"rgba(255, 103, 0, 0.2)",borderColor:"rgba(255, 103, 0, 1)",borderWidth:2,hoverBackgroundColor:"rgba(255,99,132,0.4)",hoverBorderColor:"rgba(255,99,132,1)",data:a.props.chartData404}]},chartOptions:{responsive:!0,scales:{xAxes:[{time:{unit:"day",displayFormats:{quarter:"MMM D"}},distribution:"series"}],yAxes:[{ticks:{beginAtZero:!0}}]}}},a}return Object(d.a)(t,e),Object(o.a)(t,[{key:"render",value:function(){return l.a.createElement("div",{className:"reactChart"},l.a.createElement(f.a,{data:this.state.chartData,options:this.state.chartOptions}))}}]),t}(n.Component),g={border:"3px solid #00C853"},b={border:"3px solid #D50000"},y={textAlign:"left",listStyleType:"none",margin:0,padding:"0 0 0 8px"},v=function(e){function t(){return Object(s.a)(this,t),Object(u.a)(this,Object(m.a)(t).apply(this,arguments))}return Object(d.a)(t,e),Object(o.a)(t,[{key:"render",value:function(){return l.a.createElement("div",{className:"domainCard card",style:this.props.statusCode>399?b:g},l.a.createElement("div",{className:"header"},l.a.createElement("h4",{className:"title"},l.a.createElement("a",{href:"http://"+this.props.domain,target:"_blank"},""!==this.props.statusInfo.FacilityName?this.props.statusInfo.FacilityName:this.props.domain))),l.a.createElement("div",{className:"content"+(this.props.ctAllIcons?" all-icons":"")+(this.props.ctTableFullWidth?" table-full-width":"")+(this.props.ctTableResponsive?" table-responsive":"")+(this.props.ctTableUpgrade?" table-upgrade":"")},l.a.createElement(E,{chartData:this.props.statusInfo.GraphDataOutage,chartData404:this.props.statusInfo.GraphData404}),l.a.createElement("div",{className:"footer"},null!=this.props.stats?l.a.createElement("hr",null):"",l.a.createElement("div",{className:"stats"},l.a.createElement("ul",{style:y},l.a.createElement("li",{hidden:200===this.props.statusCode},l.a.createElement("b",null,"Status:")," ",this.props.statusInfo.Status),l.a.createElement("li",null,l.a.createElement("b",null,"Response Time:")," ",this.props.statusInfo.AvgResponse.toFixed(2),l.a.createElement("i",null,"ms")),l.a.createElement("li",null,l.a.createElement("b",null,l.a.createElement("span",{style:{color:"#e83e8c"}},"Outages: ")),this.props.statusInfo.Outages),l.a.createElement("li",null,l.a.createElement("b",null,l.a.createElement("span",{style:{color:"#ff6700"}},"404: ")),this.props.statusInfo.Errors))))))}}]),t}(n.Component),C=a(47),w=a.n(C),O="http://monitoring-env.qj3cticwqw.us-east-1.elasticbeanstalk.com/api/monitorstatus",S="http://monitoring-env.qj3cticwqw.us-east-1.elasticbeanstalk.com/api/monthlymonitorstatus",j="http://monitoring-env.qj3cticwqw.us-east-1.elasticbeanstalk.com/api/404list",I="http://monitoring-env.qj3cticwqw.us-east-1.elasticbeanstalk.com/api/runCrawl";"monitor.acadiadevelopment.com"===document.location.host&&(O="http://"+document.location.host+"/api/monitorstatus",S="http://"+document.location.host+"/api/monthlymonitorstatus",j="http://"+document.location.host+"/api/404list",I="http://"+document.location.host+"/api/runCrawl");var k={crawl:{startCrawl:function(e,t,a){var n=arguments.length>3&&void 0!==arguments[3]?arguments[3]:"",l=new FormData;return l.set("domain",e),l.set("email",t),l.set("crawlType",a),l.set("searchTerm",n),w.a.post(I,l,{headers:{"Content-Type":"multipart/form-data"}}).then(function(e){return e.data}).catch(function(e){console.log(e)})}},facility:{getFacilityList:function(){return w.a.get("http://monitoring-env.qj3cticwqw.us-east-1.elasticbeanstalk.com/api/getFacilities").then(function(e){return e.data}).catch(function(e){console.log("GET Facility List ERR: ",e)})}},fof:{get404List:function(){return w.a.get(j).then(function(e){return e.data}).catch(function(e){console.log("GET 404 ERR: ",e)})}},status:{getStatusInfo:function(){return w.a.get(O).then(function(e){return e.data}).catch(function(e){console.log("GET Status ERR: ",e)})},getMonthlyStatusInfo:function(){return w.a.get(S).then(function(e){return e.data}).catch(function(e){console.log("GET MonthlyStatus ERR: ",e)})}}},D=a(298),N=a(291),F=a(289),T=a(77),A=a.n(T);var x=function(e){return l.a.createElement("div",{className:"fof"},l.a.createElement("h4",null,e.facility),l.a.createElement("table",{className:"fofTable",style:{borderCollapse:"collapse"}},l.a.createElement("thead",null,l.a.createElement("tr",null,l.a.createElement("th",{component:"th"}),l.a.createElement("td",null,"404 Link"),l.a.createElement("td",null,"Referred From / Found On"),l.a.createElement("td",null,"Time"))),l.a.createElement("tbody",null,Object.keys(e.data).map(function(t,a){return l.a.createElement("tr",{key:t},l.a.createElement("th",null,t),l.a.createElement("td",null,l.a.createElement("a",{href:e.data[a].Page.String,target:"_blank"},e.data[a].Page.String)),l.a.createElement("td",null,l.a.createElement("a",{href:e.data[a].Referer.String,target:"_blank"},e.data[a].Referer.String)),l.a.createElement("td",null,A()(e.data[a].TimeStamp).format("lll")))}))))},q=function(e){function t(){var e,a;Object(s.a)(this,t);for(var n=arguments.length,l=new Array(n),r=0;r<n;r++)l[r]=arguments[r];return(a=Object(u.a)(this,(e=Object(m.a)(t)).call.apply(e,[this].concat(l)))).state={fofList:{}},a.retrieve=function(){k.fof.get404List().then(function(e){a.setState({fofList:e})})},a}return Object(d.a)(t,e),Object(o.a)(t,[{key:"componentWillMount",value:function(){this.retrieve()}},{key:"render",value:function(){var e=this.props.selected,t=this.state.fofList;return console.log("Selected: ",e),console.log(t),this.props.selected.length>0?l.a.createElement("div",{className:"content"},l.a.createElement("h3",null,"404's"),Object.keys(t).filter(function(a){return e.includes(t[a][0].FacilityName.String)}).map(function(e,a){return l.a.createElement(x,{facility:t[e][0].FacilityName.String,domain:e,key:a,data:t[e]})})):l.a.createElement("div",{className:"content"},l.a.createElement("h3",null,"404's"),t?Object.keys(t).map(function(e,a){return l.a.createElement(x,{facility:t[e][0].FacilityName.String,domain:e,key:a,data:t[e]})}):l.a.createElement("p",null,"Could Not Get Data"))}}]),t}(n.Component);function R(e){return l.a.createElement(F.a,{component:"div",style:{padding:24}},e.children)}var M=function(e){function t(){var e,a;Object(s.a)(this,t);for(var n=arguments.length,l=new Array(n),r=0;r<n;r++)l[r]=arguments[r];return(a=Object(u.a)(this,(e=Object(m.a)(t)).call.apply(e,[this].concat(l)))).state={domainObj:{},monthlyDomainObj:{},showMonthly:!1,lastUpdate:" ",value:0},a.retrieve=function(){k.status.getStatusInfo().then(function(e){a.setState({domainObj:e})}),k.status.getMonthlyStatusInfo().then(function(e){a.setState({monthlyDomainObj:e})})},a.handleChange=function(e,t){a.setState({value:t})},a}return Object(d.a)(t,e),Object(o.a)(t,[{key:"componentWillMount",value:function(){this.retrieve()}},{key:"render",value:function(){var e=this.props.selected,t=this.state,a=t.value,n=t.domainObj,r=t.monthlyDomainObj;return console.log(n),e.length>0?l.a.createElement("div",{className:"content"},l.a.createElement("h3",null,"Site Status"),l.a.createElement(D.a,{value:a,onChange:this.handleChange,centered:!0},l.a.createElement(N.a,{label:"Weekly"}),l.a.createElement(N.a,{label:"Monthly"}),l.a.createElement(N.a,{label:"404 List"}),l.a.createElement(p.a,{to:"/manual_crawl"},l.a.createElement(N.a,{label:"Start Crawl"}))),0===a&&l.a.createElement(R,null,n?Object.keys(n).map(function(t,a){if(e.includes(n[t].FacilityName))return l.a.createElement(v,{isSelected:!0,statsIcon:"fa fa-history",key:a,statusCode:n[t].Status,statusInfo:n[t],id:t,domain:t})}):l.a.createElement("p",null," Could Not Get Data ")),1===a&&l.a.createElement(R,null,r?Object.keys(r).map(function(t,a){if(e.includes(r[t].FacilityName))return l.a.createElement(v,{isSelected:!0,statsIcon:"fa fa-history",key:a,statusCode:r[t].Status,statusInfo:r[t],id:t,domain:t})}):l.a.createElement("p",null," Could Not Get Data ")),2===a&&l.a.createElement(R,null,l.a.createElement(q,{selected:e}))):l.a.createElement("div",{className:"content"},l.a.createElement("h3",null,"Site Status"),l.a.createElement(D.a,{value:a,onChange:this.handleChange,centered:!0},l.a.createElement(N.a,{label:"Weekly"}),l.a.createElement(N.a,{label:"Monthly"}),l.a.createElement(N.a,{label:"404 List"}),l.a.createElement(p.a,{to:"/manual_crawl"},l.a.createElement(N.a,{label:"Start Crawl"}))),0===a&&l.a.createElement(R,null,n?Object.keys(n).map(function(e,t){return l.a.createElement(v,{isSelected:!0,statsIcon:"fa fa-history",key:t,statusCode:n[e].Status,statusInfo:n[e],id:e,domain:e})}):l.a.createElement("p",null,"Could Not Get Data")),1===a&&l.a.createElement(R,null,r?Object.keys(r).map(function(e,t){return l.a.createElement(v,{isSelected:!0,statsIcon:"fa fa-history",key:t,statusCode:r[e].Status,statusInfo:r[e],id:e,domain:e})}):l.a.createElement("p",null,"Could Not Get Data")),2===a&&l.a.createElement(R,null,l.a.createElement(q,{selected:e})))}}]),t}(n.Component),L=a(299),G=a(296),U=a(293),P=a(297),V=a(303),W=a(302),_=(a(253),function(e){function t(e){var a;return Object(s.a)(this,t),(a=Object(u.a)(this,Object(m.a)(t).call(this,e))).startCrawl=function(){var e=a.state.crawlDomain,t=a.state.userEmail,n=a.state.crawlType,l=a.state.searchTerm;return e.replace("https://",""),e.replace("http://",""),console.log(e),e.length<2||t.length<2?alert("Domain/Email Invalid"):t.indexOf("@acadiahealthcare.com")<0?alert("Requires Acadia Email"):void k.crawl.startCrawl(e,t,n,l).then(function(e){alert(e)})},a.updateDomain=function(e){e=(e=(e=e.replace("https://","")).replace("http://","")).replace(/\/$/,""),a.setState({crawlDomain:e})},a.updateUserEmail=function(e){a.setState({userEmail:e})},a.updateSearchTerm=function(e){a.setState({searchTerm:e})},a.handleSelect=function(e){a.setState({crawlType:e})},a.state={crawlDomain:"",userEmail:"",crawlType:"404",searchTerm:""},a}return Object(d.a)(t,e),Object(o.a)(t,[{key:"render",value:function(){var e=this;return l.a.createElement("div",{className:"content"},l.a.createElement("h3",null,"Crawl Domain"),l.a.createElement("form",{className:"",noValidate:!0,autoComplete:"off"},l.a.createElement("p",null,"Trigger a crawl and receive report in your email.",l.a.createElement("br",null),"Crawl times vary by domain."),l.a.createElement(U.a,{className:""},l.a.createElement(W.a,{htmlFor:"filterCrawlType"},"Crawl Type"),l.a.createElement(P.a,{style:{width:"120px",color:"#ff9800"},value:this.state.crawlType,onChange:function(t){return e.handleSelect(t.target.value)},inputProps:{name:"crawlType",id:"filterCrawlType"}},l.a.createElement(V.a,{value:"404"},"404"),l.a.createElement(V.a,{value:"sitemap"},"Sitemap"),l.a.createElement(V.a,{value:"search"},"Search"))),l.a.createElement("br",null),"search"===this.state.crawlType&&l.a.createElement("span",null,l.a.createElement(L.a,{required:!0,id:"standard-required",label:"Search Term",defaultValue:this.state.searchTerm,onChange:function(t){return e.updateSearchTerm(t.target.value)},className:"textField",margin:"normal"}),l.a.createElement("br",null)),l.a.createElement(L.a,{required:!0,id:"standard-required",label:"Domain",defaultValue:this.state.crawlDomain,onChange:function(t){return e.updateDomain(t.target.value)},className:"textField",margin:"normal"}),l.a.createElement("br",null),l.a.createElement(L.a,{required:!0,id:"standard-required",label:"Email",defaultValue:this.state.userAgent,onChange:function(t){return e.updateUserEmail(t.target.value)},className:"textField",margin:"normal"}),l.a.createElement("br",null),l.a.createElement(G.a,{variant:"contained",color:"primary",className:"",onClick:this.startCrawl},"Start Crawl")))}}]),t}(n.Component)),B=a(79),K=a(117),J=a(118),H=a(114),Y=a.n(H),Z=a(119),$=a(295),z=a(304),Q=[];var X=function(e){function t(e){var a;return Object(s.a)(this,t),(a=Object(u.a)(this,Object(m.a)(t).call(this,e))).savedFacilities=null!=localStorage.getItem("selectedFacilities")?localStorage.getItem("selectedFacilities"):[],a.state={inputValue:"",selectedItem:a.savedFacilities.length>0?a.savedFacilities.split(","):[]},a.handleKeyDown=function(e){var t=a.state,n=t.inputValue,l=t.selectedItem;l.length&&!n.length&&"backspace"===Y()(e)&&(l=l.slice(0,l.length-1),a.setState({selectedItem:l}),a.props.onUpdate(l),localStorage.setItem("selectedFacilities",l))},a.handleInputChange=function(e){a.setState({inputValue:e.target.value})},a.handleChange=function(e){var t=a.state.selectedItem;-1===t.indexOf(e)&&(t=[].concat(Object(B.a)(t),[e])),a.setState({inputValue:"",selectedItem:t}),a.props.onUpdate(t),localStorage.setItem("selectedFacilities",t)},a.handleDelete=function(e){return function(){var t=Object(B.a)(a.state.selectedItem);t.splice(t.indexOf(e),1),a.setState({selectedItem:t}),a.props.onUpdate(t),localStorage.setItem("selectedFacilities",t)}},Q.length<1&&k.facility.getFacilityList().then(function(e){for(var t in Q=[{label:"All CTC",url:"ctc"},{label:"All Inpatient",url:"inpatient"},{label:"All Residential Dual",url:"residentialDual"},{label:"All Residential SA",url:"residentialSA"},{label:"All Specialty",url:"specialty"}],e)Q.push({label:e[t].FacilityName,url:t,type:e[t].FacilityType})}),a}return Object(d.a)(t,e),Object(o.a)(t,[{key:"render",value:function(){var e=this,t=this.state,a=t.inputValue,n=t.selectedItem;return l.a.createElement(Z.a,{inputValue:a,onChange:this.handleChange,selectedItem:n},function(t){var a=t.getInputProps,r=t.getItemProps,c=t.isOpen,i=t.inputValue,s=t.selectedItem,o=t.highlightedIndex;return l.a.createElement("div",{className:"autoComplete-container"},function(e){var t=e.InputProps,a=e.ref,n=Object(J.a)(e,["InputProps","ref"]);return l.a.createElement(L.a,Object.assign({InputProps:Object(K.a)({inputRef:a,classes:{root:"autoComplete-inputRoot"}},t)},n))}({fullWidth:!0,InputProps:a({startAdornment:n.map(function(t){return l.a.createElement(z.a,{key:t,tabIndex:-1,label:t,className:"autoComplete-chip",onDelete:e.handleDelete(t)})}),onChange:e.handleInputChange,onKeyDown:e.handleKeyDown,placeholder:"Search facility",id:"integration-downshift-multiple"})}),c?l.a.createElement($.a,{className:"autoComplete-paper",square:!0},function(e){var t=0;return Q.filter(function(a){var n=(!e||-1!==a.label.toLowerCase().indexOf(e.toLowerCase()))&&t<10;return n&&(t+=1),n})}(i).map(function(e,t){return function(e){var t=e.suggestion,a=e.index,n=e.itemProps,r=e.highlightedIndex===a,c=(e.selectedItem||"").indexOf(t.label)>-1;return l.a.createElement(V.a,Object.assign({},n,{key:t.label,selected:r,component:"div",style:{fontWeight:c?500:400}}),t.label)}({suggestion:e,index:t,itemProps:r({item:e.label}),highlightedIndex:o,selectedItem:s})})):null)})}}]),t}(l.a.Component),ee=(a(255),a(256),function(e){function t(){var e,a;Object(s.a)(this,t);for(var n=arguments.length,l=new Array(n),r=0;r<n;r++)l[r]=arguments[r];return(a=Object(u.a)(this,(e=Object(m.a)(t)).call.apply(e,[this].concat(l)))).state={SelectedFacilities:null!=localStorage.getItem("selectedFacilities")?localStorage.getItem("selectedFacilities"):[]},a.selectedUpdate=function(e){a.setState({SelectedFacilities:e}),console.log(e)},a.baseUrl="",a}return Object(d.a)(t,e),Object(o.a)(t,[{key:"render",value:function(){var e=this;return"#/manual_crawl"===window.location.hash?l.a.createElement("div",{className:"App"},l.a.createElement("header",{className:"App-header"},l.a.createElement("div",{className:"App-title"},l.a.createElement(p.a,{to:"/"},"Acadia Monitoring"))),l.a.createElement(h.a,{path:"/manual_crawl",render:function(){return l.a.createElement(_,null)}})):l.a.createElement("div",{className:"App"},l.a.createElement("header",{className:"App-header"},l.a.createElement("div",{className:"App-title"},l.a.createElement(p.a,{to:"/"},"Acadia Monitoring"))),l.a.createElement(X,{onUpdate:this.selectedUpdate}),l.a.createElement(h.a,{exact:!0,path:"/",render:function(){return l.a.createElement(M,{selected:e.state.SelectedFacilities})}}),"#/"!==window.location.hash?l.a.createElement("div",{style:{padding:"30px"}},"You're off the path...",l.a.createElement(p.a,{to:"/"},"Go Back Home")):l.a.createElement("span",null))}}]),t}(n.Component));c.a.render(l.a.createElement(i.a,null,l.a.createElement(ee,null)),document.getElementById("root"))}},[[129,1,2]]]);
//# sourceMappingURL=main.38aa2f42.chunk.js.map