webpackJsonp([1],{"1/oy":function(t,e){},"7Xr1":function(t,e){},"9M+g":function(t,e){},FMrB:function(t,e){},GfHa:function(t,e){},Id91:function(t,e){},L4LI:function(t,e){},NHnr:function(t,e,o){"use strict";Object.defineProperty(e,"__esModule",{value:!0});var n=o("7+uW"),a=(o("qb6w"),o("9M+g"),o("E8q/")),i=o("X1Qo"),s=o("uM87"),l=o("HDlv"),r=o("czjy"),c=o("CFyO"),u=o("FGX+"),p=o("WnFU"),d=o("ghhx"),m=o("UgDT"),v=o("hpTH"),b=o("f9FU"),g=o("W/FM"),h=o("qHHr"),f=o("C9Xq"),y=o("w+xg"),S=o("xGBE"),_=o("LHeC"),P=o("UePd"),w=o("KpFv"),C=o("UqMQ"),k=o("BiGh"),N=o("sGYS"),x=o("molP"),T=o("YfH7"),L=o("Hea7"),E=o("TnPr"),A=o("n5GY"),R=o("KXjV"),I=o("Dd8w"),V=o.n(I),D=o("NYxO"),O={name:"App",data:function(){return{appName:"Gwmpd",connected:!1,songPlayed:!1}},computed:V()({},Object(D.c)({getConnectionStatus:"getConnectionStatus",getCurrentSongInfos:"getCurrentSongInfos",getStatus:"getStatusInfos"}),{volumeValue:{get:function(){return this.getStatus.volume},set:function(t){var e=this;this.setVolume(t),this.axios.post("v1/setVolume",{volume:t}).then(function(o){e.setVolume(t)})}}}),methods:V()({refresh:function(){this.$auth.refresh({success:function(t){this.$auth.token(null,t.data.token),sessionStorage.token=t.data.token},error:function(t){delete sessionStorage.token}})}},Object(D.b)(["setAllStatus","setConnectionStatus","setSong","setVolume","setPlaylist","setState"]),{toggleMuteVolume:function(){var t=this;this.axios.put("v1/toggleMuteVolume").then(function(e){t.setVolume(e.data.volume)})},forwardSong:function(){this.axios.get("v1/nextSong")},pauseSong:function(){var t=this;"pause"!==this.getStatus.state&&this.axios.get("v1/pauseSong").then(function(e){t.setState("pause"),t.songPlayed=!1})},playSong:function(){var t=this;"play"!==this.getStatus.state&&this.axios.get("v1/playSong").then(function(e){t.setState("play"),t.songPlayed=!0})},previousSong:function(){this.axios.get("v1/previousSong")},stopSong:function(){var t=this;"stop"!==this.getStatus.state&&this.axios.get("v1/stopSong").then(function(e){t.setState("stop"),t.songPlayed=!1})}}),mounted:function(){var t=this;this.$auth.watch.authenticated&&this.refresh(),this.$refreshTokenInterval=setInterval(function(){t.refresh()},55e3),this.$refreshMpdDataInterval=setInterval(function(){t.axios.get("v1/statusMPD").then(function(e){t.setAllStatus(e.data),t.setConnectionStatus(!0),"play"===e.data.state?(t.axios.get("v1/currentSong").then(function(e){t.getCurrentSongInfos.file!==e.data.file&&t.axios.get("v1/currentPlaylist").then(function(e){t.setPlaylist(e.data)}),t.setSong(e.data)}),t.songPlayed=!0):t.songPlayed=!1},function(e){t.setConnectionStatus(!1)})},1e3)}},F={render:function(){var t=this,e=t.$createElement,o=t._self._c||e;return o("div",{attrs:{id:"app"}},[!0===t.getConnectionStatus?o("div",{staticClass:"myNavBar"},[o("b-navbar",{staticClass:"navBar",attrs:{toggleable:"md",type:"dark"}},[o("b-navbar-toggle",{attrs:{target:"nav_collapse"}}),t._v(" "),o("b-navbar-brand",[t._v(t._s(t.appName))]),t._v(" "),o("b-collapse",{attrs:{"is-nav":"",id:"nav_collapse"}},[o("b-navbar-nav",[o("b-nav-item",{attrs:{to:{name:"QueueView",params:{}}}},[t._v("Queue")]),t._v(" "),o("b-nav-item",{attrs:{to:{name:"PlaylistView",params:{}}}},[t._v("Playlists")]),t._v(" "),o("b-nav-item",{attrs:{to:{name:"ConfigView",params:{}},disabled:""}},[t._v("Config")]),t._v(" "),o("b-nav-item",{attrs:{to:{name:"AboutView",params:{}}}},[t._v("About")])],1)],1)],1),t._v(" "),o("b-container",[o("b-row",{staticClass:"positionning text-center"},[o("b-col",[o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-skip_previous iconSize",attrs:{title:"Previous song"},on:{click:t.previousSong}})],1),t._v(" "),o("b-col",[o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-stop iconSize",attrs:{title:"Stop song"},on:{click:t.stopSong}})],1),t._v(" "),o("b-col",[t.songPlayed?t._e():o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-play_arrow iconSize",attrs:{title:"Play song"},on:{click:t.playSong}}),t._v(" "),t.songPlayed?o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-pause iconSize",attrs:{title:"Pause song"},on:{click:t.pauseSong}}):t._e()],1),t._v(" "),o("b-col",[o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-skip_next iconSize",attrs:{title:"Next song"},on:{click:t.forwardSong}})],1)],1),t._v(" "),o("b-row",{staticClass:"sound text-center"},[o("b-col",{attrs:{cols:"3"}},[0==t.volumeValue?o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-volume_off iconSize",attrs:{title:"Unmute sound"},on:{click:t.toggleMuteVolume}}):t._e(),t._v(" "),t.volumeValue>0&&t.volumeValue<30?o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-volume_mute iconSize",attrs:{title:"Mute sound"},on:{click:t.toggleMuteVolume}}):t._e(),t._v(" "),t.volumeValue>29&&t.volumeValue<60?o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-volume_down iconSize",attrs:{title:"Mute sound"},on:{click:t.toggleMuteVolume}}):t._e(),t._v(" "),t.volumeValue>59?o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-volume_up iconSize",attrs:{title:"Mute sound"},on:{click:t.toggleMuteVolume}}):t._e()],1),t._v(" "),o("b-col",{staticClass:"slidecontainer",attrs:{cols:"9"}},[o("b-form-input",{staticClass:"slider",attrs:{id:"volumeSlider","b-tooltip.hover.top":"",title:t.volumeValue,type:"range",min:"0",max:"100",step:5},model:{value:t.volumeValue,callback:function(e){t.volumeValue=t._n(e)},expression:"volumeValue"}})],1)],1)],1),t._v(" "),""!==t.getStatus.error?o("div",[t._v("\n      Error: "+t._s(t.getStatus.error)+"\n    ")]):t._e()],1):o("div",[o("b-alert",{attrs:{show:"",variant:"warning"}},[o("strong",[t._v(t._s(t.appName)+" is disconnected !")]),o("br"),t._v(" "),o("router-link",{staticClass:"alert-link",attrs:{to:"Login"}},[t._v("Please sign in")])],1)],1),t._v(" "),o("b-container",{attrs:{fluid:""}},[o("router-view",{attrs:{name:"SideBar"}}),t._v(" "),o("router-view")],1)],1)},staticRenderFns:[]};var M=o("VU/8")(O,F,!1,function(t){o("FMrB")},null,null).exports,U=o("/ocq"),B={name:"Login",data:function(){return{data:{rememberMe:!1,fetchUser:!1},user:"",password:"",seePassword:"password",url:"",visibilityBool:!1,visibilityIcon:"icon-visibility"}},methods:{login:function(){this.axios.defaults.baseURL=this.url;var t=this.$auth.redirect();this.$auth.login({data:{username:this.user,password:this.password},redirect:{name:t?t.from.name:"QueueView"},rememberMe:this.data.rememberMe,fetchUser:this.data.fetchUser,success:function(t){this.$interval=5,this.$auth.token(null,t.data.token),sessionStorage.token=t.data.token,localStorage.url=this.url},error:function(){delete sessionStorage.token}})},toggleVisibility:function(){!0===this.visibilityBool?(this.visibilityBool=!1,this.visibilityIcon="icon-visibility",this.seePassword="password"):(this.visibilityBool=!0,this.visibilityIcon="icon-visibility_off",this.seePassword="text")}},mounted:function(){console.log("mounted Login.vue"),"undefined"!==localStorage.url&&(this.url=localStorage.url,this.axios.defaults.baseURL=this.url)}},z={render:function(){var t=this,e=t.$createElement,o=t._self._c||e;return o("div",{},[o("b-container",{attrs:{fluid:""}},[o("b-form",{on:{submit:function(e){return e.preventDefault(),t.login(e)}}},[o("b-form-group",[o("b-input-group",{attrs:{prepend:"API's url"}},[o("b-form-input",{attrs:{placeholder:"https://your_domain.com:8060"},model:{value:t.url,callback:function(e){t.url=e},expression:"url"}})],1)],1),t._v(" "),o("b-form-group",[o("b-input-group",{attrs:{prepend:"Login"}},[o("b-form-input",{model:{value:t.user,callback:function(e){t.user=e},expression:"user"}})],1)],1),t._v(" "),o("b-form-group",[o("b-input-group",{attrs:{prepend:"Password"}},[o("b-form-input",{attrs:{type:t.seePassword},model:{value:t.password,callback:function(e){t.password=e},expression:"password"}}),t._v(" "),o("b-input-group-append",[o("b-button",{class:t.visibilityIcon,on:{click:t.toggleVisibility}})],1)],1)],1),t._v(" "),o("b-button",{staticClass:"submitButton",attrs:{type:"submit",size:"sm",variant:"primary"}},[t._v("Sign in")])],1)],1)],1)},staticRenderFns:[]};var $=o("VU/8")(B,z,!1,function(t){o("7Xr1")},null,null).exports,j={name:"QueueView",data:function(){return{connected:!1,fields:[{key:"File",label:"Filename"},"Duration",{key:"buttonPlayMusic",label:""}]}},computed:V()({},Object(D.c)({currentPlaylist:"getCurrentPlaylist",currentSong:"getCurrentSongInfos",getConnectionStatus:"getConnectionStatus",status:"getStatusInfos"}),{getMusicDuration:function(){return this.convertSecondsToString(this.status.duration)},getMusicElapsed:function(){return this.convertSecondsToString(this.status.elapsed)},musicValue:{get:function(){return this.status.elapsed},set:function(t){var e=this;this.setPosition(t),this.axios.post("v1/setPositionTimeInCurrentSong",{position:t}).then(function(o){e.setPosition(t)})}}}),methods:V()({},Object(D.b)(["setPlaylist","setPosition","setState","setID"]),{convertSecondsToString:function(t){var e=~~(t/3600),o=~~(t%3600/60),n=~~t%60,a="";return e>0&&(a+=e+":"+(o<10?"0":"")),a+=o+":"+(n<10?"0":""),a+=""+n},playSong:function(t,e){var o=this;this.axios.get("v1/playSong",{params:{pos:e}}).then(function(e){o.setID(t),!0!==o.songPlayed&&(o.setState("play"),o.songPlayed=!0)})},moveBackwardsInTime:function(){var t=this,e=this.status.elapsed-10;e<=0?this.axios.post("v1/setPositionTimeInCurrentSong",{position:0}).then(function(e){t.setPosition(0)}):this.axios.post("v1/setPositionTimeInCurrentSong",{position:e}).then(function(o){t.setPosition(e)})},moveForwardInTime:function(){var t=this,e=this.status.elapsed+10;e>this.status.duration?this.axios.get("v1/nextSong"):this.axios.post("v1/setPositionTimeInCurrentSong",{position:e}).then(function(o){t.setPosition(e)})}}),mounted:function(){var t=this;this.axios.get("v1/currentPlaylist").then(function(e){t.setPlaylist(e.data)})}},G={render:function(){var t=this,e=t.$createElement,o=t._self._c||e;return!0===t.getConnectionStatus?o("div",{},[o("b-container",[""===t.currentSong.Title&&""===t.currentSong.Album&&""===t.currentSong.Artist?o("div",{staticClass:"text-center"},[o("h5",[t._v(" ")]),t._v(" "),o("h2",{staticClass:"truncateLongText"},[t._v(t._s(t.currentSong.file))]),t._v(" "),o("h5",[t._v(" ")])]):o("div",{staticClass:"text-center"},[o("h5",{staticClass:"truncateLongText"},[""!==t.currentSong.Album?o("span",[t._v(t._s(t.currentSong.Album))]):o("span",[t._v(" ")])]),t._v(" "),o("h2",{staticClass:"truncateLongText"},[""!==t.currentSong.Title?o("span",[t._v(t._s(t.currentSong.Title))]):o("span",[t._v(" ")])]),t._v(" "),o("h5",{staticClass:"truncateLongText"},[""!==t.currentSong.Artist?o("span",[t._v(t._s(t.currentSong.Artist))]):o("span",[t._v(" ")])])])]),t._v(" "),0!==t.status.duration?o("div",{staticClass:"musicSliderAndNumber"},[o("b-container",{attrs:{fluid:""}},[o("b-row",{staticClass:"currentSongState"},[o("b-col",{staticClass:"text-left"},[o("strong",[t._v(t._s(t.getMusicElapsed))])]),t._v(" "),o("b-col",{staticClass:"text-right"},[o("strong",[t._v(t._s(t.getMusicDuration))])])],1)],1),t._v(" "),o("input",{directives:[{name:"model",rawName:"v-model.number",value:t.musicValue,expression:"musicValue",modifiers:{number:!0}}],staticClass:"musicSlider",attrs:{type:"range","b-tooltip.hover":"",title:t.getMusicElapsed,min:"0",max:t.status.duration},domProps:{value:t.musicValue},on:{__r:function(e){t.musicValue=t._n(e.target.value)},blur:function(e){t.$forceUpdate()}}}),o("br")],1):t._e(),t._v(" "),o("div",[o("b-table",{attrs:{stacked:"md",striped:"",hover:"",items:t.currentPlaylist,fields:t.fields},scopedSlots:t._u([{key:"File",fn:function(e){return[o("span",{staticClass:"toLongFilenameSong"},[t._v(t._s(e.item.File))])]}},{key:"buttonPlayMusic",fn:function(e){return[t.currentSong.Id!==e.item.ID?o("span",{staticClass:"buttonAlignRight"},[o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-play_arrow",attrs:{title:"Play song"},on:{click:function(o){t.playSong(e.item.ID,e.item.Pos)}}})],1):t._e()]}}])})],1),t._v(" "),o("div",{},[o("router-view",{attrs:{name:"SideBar"}})],1)],1):t._e()},staticRenderFns:[]};var q=o("VU/8")(j,G,!1,function(t){o("OqfE")},"data-v-0fe2e913",null).exports,H={name:"SideBar",data:function(){return{selected:[]}},computed:V()({},Object(D.c)({getConnectionStatus:"getConnectionStatus",status:"getStatusInfos"})),methods:V()({},Object(D.b)(["setConsume","setRandom","setRepeat","setSingle","setAllStatus","setPlaylist","setRepeat"]),{clearQueue:function(){var t=this;this.axios.get("v1/clearCurrentPlaylist").then(function(e){t.setPlaylist({})})},shuffle:function(){var t=this;this.axios.get("v1/shuffle").then(function(e){t.axios.get("v1/currentPlaylist").then(function(e){t.setPlaylist(e.data)})})},toggleConsume:function(){var t=this;this.axios.put("v1/toggleConsume").then(function(e){t.setConsume(e.data.consume)})},toggleRandom:function(){var t=this;this.axios.put("v1/toggleRandom").then(function(e){t.setRandom(e.data.random)})},toggleRepeat:function(){var t=this;this.axios.put("v1/toggleRepeat").then(function(e){t.setRepeat(e.data.repeat)})},updateDB:function(){this.axios.get("v1/updateDB")}}),mounted:function(){var t=this;this.axios.get("v1/statusMPD").then(function(e){t.setAllStatus(e.data),!0===e.data.consume&&t.selected.push("consume"),!0===e.data.random&&t.selected.push("random"),!0===e.data.repeat&&t.selected.push("repeat")})}},Q={render:function(){var t=this,e=t.$createElement,o=t._self._c||e;return!0===t.getConnectionStatus?o("div",{},[o("div",{staticClass:"text-center positionning"},[o("b-row",[o("b-col",[o("b-form-group",[o("b-form-checkbox-group",{attrs:{buttons:"","button-variant":"primary"},model:{value:t.selected,callback:function(e){t.selected=e},expression:"selected"}},[o("b-form-checkbox",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-whatshot iconSize",attrs:{title:"Consume music (2-state button)",value:"consume"},on:{change:t.toggleConsume}}),t._v(" "),o("b-form-checkbox",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-shuffle iconSize",attrs:{title:"Random without changing list order (2-state button)",value:"random"},on:{change:t.toggleRandom}}),t._v(" "),o("b-form-checkbox",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-repeat iconSize",attrs:{title:"Loop playlist (2-state button)",value:"repeat"},on:{change:t.toggleRepeat}})],1)],1)],1),t._v(" "),o("b-col",[o("b-button-group",[o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-shuffle iconSize",attrs:{title:"Random list order",variant:"success"},on:{click:t.shuffle}}),t._v(" "),o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-autorenew iconSize",attrs:{title:"Update music database",variant:"success"},on:{click:t.updateDB}}),t._v(" "),o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-clear_all iconSize",attrs:{title:"Clear queue",variant:"success"},on:{click:t.clearQueue}})],1)],1)],1)],1)]):t._e()},staticRenderFns:[]};var Y=o("VU/8")(H,Q,!1,function(t){o("L4LI")},"data-v-18e4df96",null).exports,W={name:"PlaylistView",data:function(){return{fields:[{key:"name",sortable:!0,label:"Playlist's name"},{key:"removePlaylist",label:""}],newPlaylist:"",sortBy:"name",sortDesc:!1}},computed:V()({},Object(D.c)({allPlaylists:"getAllPlaylists",getConnectionStatus:"getConnectionStatus"})),methods:V()({},Object(D.b)(["setState","setAllPlaylists"]),{addNewPlaylist:function(){var t=this;""!==this.newPlaylist&&this.axios.post("v1/savePlaylist",{playlistName:this.newPlaylist}).then(function(e){t.newPlaylist="",t.loadAllPlaylists()})},clearAndLoadPlaylist:function(t){var e=this;this.axios.get("v1/clearCurrentPlaylist").then(function(o){e.loadPlaylist(t)})},editPlaylist:function(t){this.$router.push({name:"EditPlaylistView",params:{playlistName:t}})},loadAllPlaylists:function(){var t=this;this.axios.get("v1/allPlaylists").then(function(e){t.setAllPlaylists(e.data)})},loadPlaylist:function(t){var e=this;this.axios.get("v1/loadPlaylist",{params:{name:t}}).then(function(t){e.axios.get("v1/playSong").then(function(t){e.setState("play"),e.songPlayed=!0})})},removePlaylist:function(t){var e=this;this.axios.post("v1/removePlaylist",{playlistName:t}).then(function(t){e.loadAllPlaylists()})}}),mounted:function(){this.loadAllPlaylists()}},X={render:function(){var t=this,e=t.$createElement,o=t._self._c||e;return!0===t.getConnectionStatus?o("div",{staticClass:"positionning"},[o("b-container",{staticClass:"addingPlaylist"},[o("b-input-group",{attrs:{prepend:"New name"}},[o("b-form-input",{nativeOn:{keyup:function(e){return"button"in e||!t._k(e.keyCode,"enter",13,e.key,"Enter")?t.addNewPlaylist(e):null}},model:{value:t.newPlaylist,callback:function(e){t.newPlaylist=e},expression:"newPlaylist"}}),t._v(" "),o("b-input-group-append",[o("b-button",{attrs:{variant:"info"},on:{click:t.addNewPlaylist}},[t._v("Add")])],1)],1)],1),t._v(" "),o("b-table",{attrs:{stacked:"md",striped:"",hover:"","sort-by":t.sortBy,"sort-desc":t.sortDesc,items:t.allPlaylists,fields:t.fields},on:{"update:sortBy":function(e){t.sortBy=e},"update:sortDesc":function(e){t.sortDesc=e}},scopedSlots:t._u([{key:"removePlaylist",fn:function(e){return[o("div",{staticClass:"alignButtonInTable"},[o("b-button",{staticClass:"icon-mode_edit",on:{click:function(o){t.editPlaylist(e.item.name)}}}),t._v(" "),o("b-button",{staticClass:"icon-add",on:{click:function(o){t.clearAndLoadPlaylist(e.item.name)}}}),t._v(" "),o("b-button",{staticClass:"icon-queue_music",on:{click:function(o){t.loadPlaylist(e.item.name)}}}),t._v(" "),o("b-button",{staticClass:"icon-delete",on:{click:function(o){t.removePlaylist(e.item.name)}}})],1)]}}])})],1):t._e()},staticRenderFns:[]};var J=o("VU/8")(W,X,!1,function(t){o("ez39")},"data-v-df0be23c",null).exports,K=o("DAYN"),Z={components:{draggable:o.n(K).a},name:"EditPlaylistView",data:function(){return{available:{directories:[],songs:[]},directoriesFields:[{key:"nameFolder",label:"Folder"},{key:"button",label:""}],isDragging:!1,isReadOnly:!0,location:"",newPlaylistName:"",playlist:[],playlistName:"",songsFields:[{key:"Song",label:"Filename"},"Duration",{key:"addSong",label:""}],songOnMove:{},webradioURL:""}},computed:V()({},Object(D.c)({getConnectionStatus:"getConnectionStatus"})),methods:{onMove:function(t){t.relatedContext;var e=t.draggedContext;this.songOnMove=e},addSongToPlaylist:function(t){var e=this;""!==this.location&&(t=this.location+"/"+t),this.axios.post("v1/addSongToPlaylist",{songFilename:t,playlistName:this.playlistName}).then(function(t){e.getPlaylist()})},addURL:function(){var t=this;""!==this.webradioURL&&this.axios.post("v1/addSongToPlaylist",{songFilename:this.webradioURL,playlistName:this.playlistName}).then(function(e){t.getPlaylist(),t.webradioURL=""})},canEditPlaylistName:function(){this.isReadOnly=!1},cancelEditPlaylistName:function(){this.isReadOnly=!0,this.newPlaylistName=this.playlistName},checkFilesList:function(t){""===this.location?this.location=t:this.location+="/"+t,this.getFilesList()},clearPlaylist:function(){var t=this;this.axios.post("v1/clearPlaylist",{playlistName:this.playlistName}).then(function(e){t.playlist={}})},convertSecondsToString:function(t){var e=~~(t/3600),o=~~(t%3600/60),n=~~t%60,a="";return e>0&&(a+=e+":"+(o<10?"0":"")),a+=o+":"+(n<10?"0":""),a+=""+n},getFilesList:function(){var t=this;this.axios.get("v1/filesList",{params:{location:this.location}}).then(function(e){t.available=e.data;var o,n=[];for(o=0;o<t.available.directories.length;o++)n.push({Name:t.available.directories[o]});for(t.available.directories=n,n=[],o=0;o<t.available.songs.length;o++)t.available.songs[o].Duration=t.convertSecondsToString(t.available.songs[o].Duration)})},getMusicName:function(t){var e=t.split("/");return e[e.length-1]},getPlaylist:function(){var t=this;this.axios.get("v1/playlistSongsList",{params:{playlistName:this.playlistName}}).then(function(e){var o;for(t.playlist=e.data,o=0;o<t.playlist.length;o++)t.playlist[o].Order=o,t.playlist[o].Duration=t.convertSecondsToString(t.playlist[o].Duration)})},pathDown:function(){if(""!==location){var t=this.location.split("/");this.location=t.slice(0,t.length-1).join("/"),this.getFilesList()}},removeSong:function(t){var e=this;this.axios.post("v1/removeSong",{playlistName:this.playlistName,pos:t}).then(function(t){e.getPlaylist()})},renamePlaylist:function(){var t=this;this.isReadOnly=!0,this.axios.post("v1/renamePlaylist",{oldName:this.playlistName,newName:this.newPlaylistName}).then(function(e){t.playlistName=e.data.newName,t.newPlaylistName=t.playlistName,t.$router.push({name:"EditPlaylistView",params:{playlistName:t.newPlaylistName}})})}},watch:{isDragging:function(t){var e=this;!1===t&&this.songOnMove.index!==this.songOnMove.futureIndex&&this.axios.post("v1/moveSong",{playlistName:this.playlistName,oldPos:this.songOnMove.index,newPos:this.songOnMove.futureIndex}).then(function(t){e.getPlaylist()})}},mounted:function(){this.playlistName=this.$route.params.playlistName,this.newPlaylistName=this.playlistName,this.getPlaylist(),this.getFilesList()}},tt={render:function(){var t=this,e=t.$createElement,o=t._self._c||e;return!0===t.getConnectionStatus?o("div",{},[o("b-container",[o("b-input-group",[o("b-form-input",{attrs:{readonly:t.isReadOnly},nativeOn:{keyup:function(e){return"button"in e||!t._k(e.keyCode,"enter",13,e.key,"Enter")?t.renamePlaylist(e):null}},model:{value:t.newPlaylistName,callback:function(e){t.newPlaylistName=e},expression:"newPlaylistName"}}),t._v(" "),t.isReadOnly?o("b-input-group-append",[o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-mode_edit",attrs:{title:"Edit playlist's name",variant:"primary"},on:{click:t.canEditPlaylistName}})],1):o("b-input-group-append",[o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-save",attrs:{title:"Save",variant:"primary"},on:{click:t.renamePlaylist}}),t._v(" "),o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-close",attrs:{title:"Cancel",variant:"primary"},on:{click:t.cancelEditPlaylistName}})],1)],1)],1),t._v(" "),o("div",{staticClass:"saveClearButton",attrs:{align:"right"}},[o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-clear_all",attrs:{size:"lg",title:"Clear playlist",variant:"danger"},on:{click:t.clearPlaylist}})],1),t._v(" "),o("b-list-group",[o("draggable",{attrs:{list:t.playlist,move:t.onMove},on:{start:function(e){t.isDragging=!0},end:function(e){t.isDragging=!1}}},t._l(t.playlist,function(e,n){return o("b-list-group-item",{key:e.Order,staticClass:"d-flex justify-content-between align-items-center"},[o("span",{staticClass:"toLongFilenameSong"},[t._v(t._s(e.File))]),t._v(" "),o("div",{staticClass:"buttonAlignRightInListGroup"},["0:00"!==e.Duration?o("b-badge",{staticClass:"addSpace",attrs:{pill:""}},[t._v(t._s(e.Duration))]):t._e(),t._v(" "),o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-delete",attrs:{title:"Remove from playlist"},on:{click:function(e){t.removeSong(n)}}})],1)])}))],1),t._v(" "),o("br"),t._v(" "),o("hr"),t._v(" "),o("b-form-group",[o("b-input-group",{attrs:{prepend:"web's URL"}},[o("b-form-input",{nativeOn:{keyup:function(e){return"button"in e||!t._k(e.keyCode,"enter",13,e.key,"Enter")?t.addURL(e):null}},model:{value:t.webradioURL,callback:function(e){t.webradioURL=e},expression:"webradioURL"}}),t._v(" "),o("b-input-group-append",[o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-add",attrs:{title:"Add a web url"},on:{click:t.addURL}})],1)],1)],1),t._v(" "),o("hr"),t._v(" "),o("b-input-group",{attrs:{prepend:"Location"}},[o("b-form-input",{attrs:{placeholder:"/",readonly:""},model:{value:t.location,callback:function(e){t.location=e},expression:"location"}}),t._v(" "),o("b-input-group-append",[""===t.location?o("b-button",{staticClass:"icon-undo",attrs:{disabled:""},on:{click:t.pathDown}}):o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-undo",attrs:{title:"Go back in the tree"},on:{click:t.pathDown}})],1)],1),t._v(" "),o("div",[t.available.directories.length>0?o("b-table",{attrs:{stacked:"md",striped:"",hover:"",items:t.available.directories,fields:t.directoriesFields},scopedSlots:t._u([{key:"nameFolder",fn:function(e){return[o("span",{staticClass:"toLongFilenameSong"},[t._v(t._s(e.item.Name))])]}},{key:"button",fn:function(e){return[o("div",{staticClass:"buttonAlignRight"},[o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-visibility",attrs:{title:"Go inside"},on:{click:function(o){t.checkFilesList(e.item.Name)}}}),t._v(" "),o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-create_new_folder",attrs:{title:"Add all the songs from the folder"},on:{click:function(o){t.addSongToPlaylist(e.item.Name)}}})],1)]}}])}):t._e()],1),t._v(" "),o("div",[t.available.songs.length>0?o("b-table",{attrs:{stacked:"md",striped:"",hover:"",items:t.available.songs,fields:t.songsFields},scopedSlots:t._u([{key:"Song",fn:function(e){return[o("span",{staticClass:"toLongFilenameSong"},[t._v(t._s(e.item.File))])]}},{key:"Duration",fn:function(e){return[o("span",[t._v(t._s(e.item.Duration))])]}},{key:"addSong",fn:function(e){return[o("b-button",{directives:[{name:"b-tooltip",rawName:"v-b-tooltip.hover.top",modifiers:{hover:!0,top:!0}}],staticClass:"icon-add buttonAlignRight",attrs:{title:"Add song"},on:{click:function(o){t.addSongToPlaylist(e.item.File)}}})]}}])}):t._e()],1)],1):t._e()},staticRenderFns:[]};var et=o("VU/8")(Z,tt,!1,function(t){o("ngK6")},"data-v-6baac48b",null).exports,ot={render:function(){var t=this,e=t.$createElement,n=t._self._c||e;return n("div",{staticClass:"text-center topSpace"},[n("b-container",{staticClass:"mx-auto backgroundAndRoundedBox",attrs:{fluid:""}},[n("h2",[t._v("Gwmpd")]),t._v(" "),n("span",{staticClass:"text-primary"},[t._v("Freyermuth Julien")]),n("br"),t._v("\n    Copyright © 2018"),n("br"),t._v(" "),n("a",{attrs:{href:"https://en.wikipedia.org/wiki/BSD_licenses",target:"_blank"}},[n("b-img",{staticClass:"licenseImg",attrs:{src:o("rO82"),alt:"Open Source logo"}})],1),t._v(" "),n("div",{staticClass:"donate"},[t._v("\n      If you like my work, you can donate to support it. I will be even more motivated to add new features and improve the documentation"),n("br"),t._v(" "),n("form",{staticClass:"topSpace",attrs:{action:"https://www.paypal.com/cgi-bin/webscr",method:"post",target:"_top"}},[n("input",{attrs:{type:"hidden",name:"cmd",value:"_s-xclick"}}),t._v(" "),n("input",{attrs:{type:"hidden",name:"hosted_button_id",value:"8SJ9D6CN28PSE"}}),t._v(" "),n("input",{attrs:{type:"image",src:"https://www.paypalobjects.com/fr_FR/FR/i/btn/btn_donateCC_LG.gif",border:"0",name:"submit",alt:"PayPal, le réflexe sécurité pour payer en ligne"}}),t._v(" "),n("img",{attrs:{alt:"",border:"0",src:"https://www.paypalobjects.com/fr_FR/i/scr/pixel.gif",width:"1",height:"1"}})])])])],1)},staticRenderFns:[]};var nt=o("VU/8")({name:"about"},ot,!1,function(t){o("WLUp")},"data-v-3d66db00",null).exports;n.a.use(U.a);var at=new U.a({mode:"history",routes:[{path:"/queue",name:"QueueView",meta:{auth:!0},components:{default:q,SideBar:Y}},{path:"/login",name:"Login",meta:{auth:!1},component:$},{path:"/about",name:"AboutView",components:{default:nt}},{path:"/playlist",name:"PlaylistView",meta:{auth:!0},components:{default:J}},{path:"/editPlaylist/:playlistName",name:"EditPlaylistView",meta:{auth:!0},components:{default:et}},{path:"*",redirect:"/queue"}]}),it=o("mtWM"),st=o.n(it),lt=o("Rf8U"),rt=o.n(lt),ct=o("MLZB"),ut=o.n(ct);n.a.use(D.a);var pt={state:{playlist:{}},getters:{getCurrentPlaylist:function(t){return t.playlist}},mutations:{SET_PLAYLIST:function(t,e){var o,n,a,i,s,l,r,c;for(o=0;o<e.length;o++)e[o].Duration=(i=e[o].Duration,void 0,void 0,void 0,c=void 0,l=~~(i%3600/60),r=~~i%60,c="",(s=~~(i/3600))>0&&(c+=s+":"+(l<10?"0":"")),c+=l+":"+(r<10?"0":""),c+=""+r),e[o].File=(n=e[o].File,void 0,(a=n.split("/"))[a.length-1]);t.playlist=e}},actions:{setPlaylist:function(t,e){t.state;(0,t.commit)("SET_PLAYLIST",e)}}},dt={state:{status:{}},getters:{getStatusInfos:function(t){return t.status}},mutations:{SET_ALL_STATUS:function(t,e){t.status=e},SET_CONSUME:function(t,e){t.status.consume=e},SET_POSITION:function(t,e){t.status.elapsed=Math.round(100*e)/100},SET_RANDOM:function(t,e){t.status.random=e},SET_REPEAT:function(t,e){t.status.repeat=e},SET_SINGLE:function(t,e){t.status.single=e},SET_STATE:function(t,e){t.status.state=e},SET_VOLUME:function(t,e){t.status.volume=e}},actions:{setAllStatus:function(t,e){t.state;(0,t.commit)("SET_ALL_STATUS",e)},setConsume:function(t,e){t.state;(0,t.commit)("SET_CONSUME",e)},setPosition:function(t,e){t.state;(0,t.commit)("SET_POSITION",e)},setRandom:function(t,e){t.state;(0,t.commit)("SET_RANDOM",e)},setRepeat:function(t,e){t.state;(0,t.commit)("SET_REPEAT",e)},setSingle:function(t,e){t.state;(0,t.commit)("SET_SINGLE",e)},setState:function(t,e){t.state;(0,t.commit)("SET_STATE",e)},setVolume:function(t,e){t.state;(0,t.commit)("SET_VOLUME",e)}}},mt=new D.a.Store({modules:{allPlaylists:{state:{allPlaylists:[]},getters:{getAllPlaylists:function(t){return t.allPlaylists}},mutations:{SET_ALL_PLAYLISTS:function(t,e){var o,n=[];for(o=0;o<e.length;o++)n.push({name:e[o]});t.allPlaylists=n}},actions:{setAllPlaylists:function(t,e){t.state,(0,t.commit)("SET_ALL_PLAYLISTS",e)}}},currentPlaylist:pt,currentSong:{state:{song:{}},getters:{getCurrentSongInfos:function(t){return t.song}},mutations:{SET_SONG:function(t,e){t.song=e},SET_ID:function(t,e){t.song.Id=e}},actions:{setSong:function(t,e){t.state,(0,t.commit)("SET_SONG",e)},setID:function(t,e){t.state,(0,t.commit)("SET_ID",e)}}},connectionStatus:{state:{status:{}},getters:{getConnectionStatus:function(t){return t.status}},mutations:{SET_CONNECTION_STATUS:function(t,e){t.status=e}},actions:{setConnectionStatus:function(t,e){t.state,(0,t.commit)("SET_CONNECTION_STATUS",e)}}},status:dt},strict:!0});n.a.router=at,n.a.component("b-button",a.a),n.a.component("b-button-group",i.a),n.a.component("b-button-toolbar",s.a),n.a.component("b-alert",l.a),n.a.component("b-form",r.a),n.a.component("b-form-group",c.a),n.a.component("b-form-input",u.a),n.a.component("b-form-checkbox",k.a),n.a.component("b-form-checkbox-group",N.a),n.a.component("b-container",p.a),n.a.component("b-row",d.a),n.a.component("b-col",m.a),n.a.component("b-img",v.a),n.a.component("b-navbar",b.a),n.a.component("b-navbar-brand",g.a),n.a.component("b-navbar-nav",h.a),n.a.component("b-navbar-toggle",f.a),n.a.component("b-nav",y.a),n.a.component("b-nav-item",S.a),n.a.component("b-nav-form",_.a),n.a.component("b-collapse",P.a),n.a.component("b-progress-bar",w.a),n.a.component("b-table",C.a),n.a.directive("b-tooltip",x.a),n.a.component("b-input-group",T.a),n.a.component("b-input-group-append",L.a),n.a.component("b-list-group",E.a),n.a.component("b-list-group-item",A.a),n.a.component("b-badge",R.a),n.a.use(rt.a,st.a),n.a.axios.defaults.baseURL="http://localhost:8060",n.a.use(ut.a,{authRedirect:"Login",auth:o("2T3s"),http:o("E/+Z"),router:o("LFDJ"),loginData:{url:"login",method:"POST",redirect:"/queue"},refreshData:{url:"v1/refresh",enabled:!1},fetchData:{enabled:!1}}),n.a.config.productionTip=!1,new n.a({el:"#app",router:at,store:mt,components:{App:M},template:"<App/>"})},OqfE:function(t,e){},WLUp:function(t,e){},ez39:function(t,e){},ngK6:function(t,e){},qb6w:function(t,e){},rO82:function(t,e,o){t.exports=o.p+"static/img/license.f54ea6f.png"}},["NHnr"]);
//# sourceMappingURL=app.57eb23db086b03be990c.js.map