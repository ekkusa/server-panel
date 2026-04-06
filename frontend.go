package main

import "net/http"

func serveFrontend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(indexHTML))
}

const indexHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Miyoubi Panel</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
:root {
  --bg:#0a0a0a;--sb:#111;--card:#141414;--card2:#1a1a1a;
  --border:#222;--border2:#2a2a2a;--text:#e8e8e8;--muted:#666;--muted2:#444;
  --accent:#fff;--green:#22c55e;--red:#ef4444;--amber:#f59e0b;--blue:#3b82f6;
  --font:'Inter',sans-serif;--mono:'JetBrains Mono',monospace;
}
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
html,body{height:100%;background:var(--bg);color:var(--text);font-family:var(--font);overflow:hidden;font-size:14px}

/* Login */
#login-overlay{position:fixed;inset:0;z-index:1000;background:var(--bg);display:flex;align-items:center;justify-content:center}
.login-box{width:360px;background:var(--card);border:1px solid var(--border2);border-radius:10px;padding:32px}
.login-logo{display:flex;align-items:center;gap:10px;margin-bottom:28px}
.login-logo-icon{width:34px;height:34px;background:var(--card2);border:1px solid var(--border2);border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:16px}
.login-logo-text{font-size:1.05rem;font-weight:600}
.login-logo-text span{color:var(--muted);font-weight:400}
.login-title{font-size:1.1rem;font-weight:600;margin-bottom:4px}
.login-sub{font-size:0.78rem;color:var(--muted);margin-bottom:24px}
.form-group{margin-bottom:14px}
.form-label{font-size:0.72rem;font-weight:500;color:var(--muted);letter-spacing:0.04em;text-transform:uppercase;display:block;margin-bottom:6px}
.form-input{width:100%;background:var(--card2);border:1px solid var(--border2);border-radius:6px;padding:9px 12px;font-family:var(--font);font-size:0.85rem;color:var(--text);outline:none;transition:border-color .15s}
.form-input:focus{border-color:#444}
.form-input::placeholder{color:var(--muted2)}
.login-btn{width:100%;margin-top:8px;background:var(--accent);color:#000;border:none;border-radius:6px;padding:10px;font-family:var(--font);font-size:0.85rem;font-weight:600;cursor:pointer;transition:opacity .15s}
.login-btn:hover{opacity:.88}
.login-error{font-size:0.75rem;color:var(--red);margin-top:10px;display:none;text-align:center}

/* Shell */
.app{display:flex;height:100vh;overflow:hidden}

/* Sidebar */
.sidebar{width:210px;min-width:210px;background:var(--sb);border-right:1px solid var(--border);display:flex;flex-direction:column;overflow:hidden}
.sb-logo{padding:16px;display:flex;align-items:center;gap:10px;border-bottom:1px solid var(--border)}
.sb-logo-icon{width:26px;height:26px;background:var(--card2);border:1px solid var(--border2);border-radius:6px;display:flex;align-items:center;justify-content:center;font-size:12px}
.sb-logo-text{font-size:0.9rem;font-weight:600}
.sb-section{padding:14px 12px 4px;font-size:0.62rem;font-weight:500;color:var(--muted2);letter-spacing:.1em;text-transform:uppercase}
.nav-item{display:flex;align-items:center;gap:9px;padding:7px 12px;margin:1px 6px;border-radius:5px;cursor:pointer;font-size:0.8rem;font-weight:500;color:var(--muted);transition:color .12s,background .12s;text-decoration:none}
.nav-item:hover{color:var(--text);background:rgba(255,255,255,.04)}
.nav-item.active{color:var(--text);background:var(--card2)}
.nav-icon{width:14px;height:14px;flex-shrink:0;opacity:.7}
.nav-item.active .nav-icon{opacity:1}
.qa-item{display:flex;align-items:center;gap:8px;padding:6px 12px;margin:1px 6px;border-radius:5px;cursor:pointer;transition:background .12s}
.qa-item:hover{background:rgba(255,255,255,.03)}
.qa-item.active{background:var(--card2)}
.qa-dot{width:6px;height:6px;border-radius:50%;background:var(--muted2);flex-shrink:0}
.qa-dot.on{background:var(--green)}
.qa-name{font-size:0.78rem;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.sb-spacer{flex:1}
.sb-user{padding:10px 12px;border-top:1px solid var(--border);display:flex;align-items:center;gap:9px;cursor:pointer;transition:background .12s}
.sb-user:hover{background:rgba(255,255,255,.03)}
.sb-avatar{width:26px;height:26px;border-radius:50%;background:var(--card2);border:1px solid var(--border2);display:flex;align-items:center;justify-content:center;font-size:0.65rem;color:var(--muted);font-weight:600;flex-shrink:0}
.sb-user-name{font-size:0.78rem;font-weight:500}
.sb-user-role{font-size:0.65rem;color:var(--muted)}
.sb-logout{margin-left:auto;font-size:0.6rem;color:var(--muted2);padding:3px 6px;border:1px solid var(--border);border-radius:4px;transition:color .12s,border-color .12s;white-space:nowrap;flex-shrink:0}
.sb-logout:hover{color:var(--red);border-color:var(--red)}
.sb-version{font-size:0.6rem;color:var(--muted2);padding:4px 14px 8px}

/* Main */
.main{flex:1;display:flex;flex-direction:column;overflow:hidden}
.srv-header{padding:14px 22px;border-bottom:1px solid var(--border);display:flex;align-items:center;gap:14px;background:var(--sb);flex-shrink:0}
.srv-icon{width:44px;height:44px;border-radius:8px;background:var(--card2);border:1px solid var(--border2);display:flex;align-items:center;justify-content:center;font-size:18px;flex-shrink:0}
.srv-title{font-size:1.1rem;font-weight:600;letter-spacing:-.01em}
.srv-desc{font-size:0.72rem;color:var(--muted);margin-top:2px}
.srv-meta{font-size:0.65rem;color:var(--muted2);margin-top:2px}
.srv-actions{margin-left:auto;display:flex;gap:8px;align-items:center;flex-shrink:0}
.act-btn{font-size:0.75rem;font-weight:500;padding:6px 14px;border-radius:5px;border:1px solid var(--border2);cursor:pointer;background:var(--card2);color:var(--text);display:flex;align-items:center;gap:6px;transition:background .12s,border-color .12s,color .12s}
.act-btn:disabled{opacity:.3;cursor:not-allowed}
.act-btn:hover:not(:disabled){background:var(--card);border-color:#555}
.act-btn-stop:hover:not(:disabled){border-color:var(--red);color:var(--red)}
.act-btn-restart:hover:not(:disabled){border-color:var(--amber);color:var(--amber)}
.act-btn-start:hover:not(:disabled){border-color:var(--green);color:var(--green)}

/* Stats */
.content-wrap{flex:1;display:flex;flex-direction:column;overflow:hidden}
.content-inner{padding:16px 22px 0;flex-shrink:0}
.stats-row{display:grid;grid-template-columns:1fr 1fr 1fr 1.1fr;gap:10px;margin-bottom:14px}
.stat-card{background:var(--card);border:1px solid var(--border);border-radius:8px;padding:14px}
.stat-card-header{display:flex;justify-content:space-between;align-items:flex-start;margin-bottom:10px}
.stat-card-label{font-size:0.65rem;font-weight:500;color:var(--muted);letter-spacing:.06em;text-transform:uppercase}
.stat-card-sub{font-size:0.6rem;color:var(--muted2);margin-top:2px}
.stat-card-icon{width:28px;height:28px;border-radius:6px;background:var(--card2);border:1px solid var(--border2);display:flex;align-items:center;justify-content:center;font-size:13px;flex-shrink:0}
.waveform{display:flex;align-items:flex-end;gap:3px;height:36px;margin-bottom:8px}
.w-bar{width:4px;border-radius:2px;background:var(--border2);height:6px;transition:background .4s}
.waveform.running .w-bar{background:var(--green);animation:wave 1.1s ease-in-out infinite}
.waveform.running .w-bar:nth-child(1){animation-delay:.00s}
.waveform.running .w-bar:nth-child(2){animation-delay:.12s}
.waveform.running .w-bar:nth-child(3){animation-delay:.24s}
.waveform.running .w-bar:nth-child(4){animation-delay:.36s}
.waveform.running .w-bar:nth-child(5){animation-delay:.48s}
.waveform.running .w-bar:nth-child(6){animation-delay:.60s}
@keyframes wave{0%,100%{height:5px}50%{height:30px}}
.stat-status-text{font-size:1rem;font-weight:600;color:var(--muted)}
.stat-status-text.running{color:var(--green)}
.stat-status-text.stopped{color:var(--red)}
.stat-status-sub{font-size:0.65rem;color:var(--muted);margin-top:2px}
.conn-address{background:var(--card2);border:1px solid var(--border2);border-radius:5px;padding:8px 10px;font-family:var(--mono);font-size:0.72rem;display:flex;align-items:center;justify-content:space-between;cursor:pointer;transition:border-color .15s;margin-top:4px}
.conn-address:hover{border-color:#444}
.copy-label{font-size:0.6rem;color:var(--muted2);font-family:var(--font)}
.conn-hint{font-size:0.62rem;color:var(--muted2);margin-top:6px}
.info-row{display:flex;justify-content:space-between;align-items:center;padding:4px 0;border-bottom:1px solid var(--border)}
.info-row:last-child{border-bottom:none}
.info-key{font-size:0.65rem;color:var(--muted)}
.info-val{font-size:0.65rem;color:var(--text);font-family:var(--mono);max-width:55%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.perf-item{margin-bottom:8px}
.perf-item:last-child{margin-bottom:0}
.perf-header{display:flex;justify-content:space-between;margin-bottom:4px}
.perf-label{font-size:0.62rem;font-weight:500;color:var(--muted);text-transform:uppercase;letter-spacing:.06em}
.perf-val{font-size:0.62rem;font-weight:600;font-family:var(--mono)}
.perf-bar-track{height:4px;background:var(--border2);border-radius:2px;overflow:hidden}
.perf-bar-fill{height:100%;border-radius:2px;background:#555;transition:width .6s ease,background .4s;width:0%}
.perf-bar-fill.warn{background:var(--amber)}
.perf-bar-fill.danger{background:var(--red)}

/* Tabs */
.tabs-bar{display:flex;border-bottom:1px solid var(--border);padding:0 22px;background:var(--sb);flex-shrink:0}
.tab-btn{font-size:0.78rem;font-weight:500;color:var(--muted);padding:9px 14px;border:none;background:transparent;cursor:pointer;border-bottom:2px solid transparent;margin-bottom:-1px;transition:color .12s,border-color .12s}
.tab-btn:hover{color:var(--text)}
.tab-btn.active{color:var(--text);border-bottom-color:var(--text)}
.tab-btn:disabled{opacity:.25;cursor:not-allowed}

/* Tab panels */
.tab-content{flex:1;display:flex;flex-direction:column;overflow:hidden}
.tab-panel{display:none;flex:1;flex-direction:column;overflow:hidden;padding:14px 22px 10px}
.tab-panel.active{display:flex}

/* Console */
.console-toolbar{display:flex;align-items:center;gap:8px;margin-bottom:8px;flex-shrink:0}
.console-title{font-size:0.82rem;font-weight:500;display:flex;align-items:center;gap:7px}
.status-badge{font-size:0.58rem;font-weight:500;padding:2px 7px;border-radius:10px;background:var(--card2);color:var(--muted);border:1px solid var(--border2)}
.status-badge.running{background:rgba(34,197,94,.1);color:var(--green);border-color:rgba(34,197,94,.2)}
.status-badge.stopped{background:rgba(239,68,68,.1);color:var(--red);border-color:rgba(239,68,68,.2)}
.console-actions{margin-left:auto;display:flex;gap:5px;align-items:center}
.con-btn{background:var(--card2);border:1px solid var(--border2);color:var(--muted);padding:4px 8px;border-radius:4px;cursor:pointer;font-size:0.65rem;transition:color .12s,border-color .12s}
.con-btn:hover{color:var(--text);border-color:#444}

#console{flex:1;background:#0d0d0d;border:1px solid var(--border);border-radius:6px;overflow-y:auto;padding:10px 12px;font-family:var(--mono);font-size:0.72rem;line-height:1.75;color:#888}
#console::-webkit-scrollbar{width:4px}
#console::-webkit-scrollbar-track{background:transparent}
#console::-webkit-scrollbar-thumb{background:var(--border2);border-radius:2px}
.log-line{word-break:break-all}
.log-warn{color:var(--amber)}
.log-error{color:var(--red)}
.log-info{color:#ccc}
.console-input-row{display:flex;margin-top:7px;flex-shrink:0;border:1px solid var(--border2);border-radius:5px;overflow:hidden;transition:border-color .15s}
.console-input-row:focus-within{border-color:#444}
.input-prompt{background:var(--card2);border-right:1px solid var(--border2);padding:0 10px;display:flex;align-items:center;color:var(--muted);font-family:var(--mono);font-size:0.78rem}
#cmd-input{flex:1;background:var(--card);border:none;outline:none;color:var(--text);padding:8px 10px;font-family:var(--mono);font-size:0.75rem}#cmd-input:disabled{color:var(--muted2);cursor:not-allowed}.console-offline-msg{display:none;align-items:center;justify-content:center;gap:8px;padding:8px 12px;margin-top:7px;background:var(--card2);border:1px solid var(--border2);border-radius:5px;font-size:0.72rem;color:var(--muted);flex-shrink:0}.console-offline-msg.visible{display:flex}
#cmd-input::placeholder{color:var(--muted2)}
.cmd-send{background:var(--card2);border:none;border-left:1px solid var(--border2);color:var(--muted);padding:0 12px;cursor:pointer;display:flex;align-items:center;transition:background .12s,color .12s}
.cmd-send:hover{background:var(--card);color:var(--text)}
.console-footer{display:flex;align-items:center;gap:14px;padding:5px 2px;font-size:0.65rem;color:var(--muted);flex-shrink:0}
.console-footer label{display:flex;align-items:center;gap:5px;cursor:pointer;user-select:none}
.console-footer input[type=checkbox]{accent-color:#555;cursor:pointer}
.console-footer select{background:var(--card2);border:1px solid var(--border2);color:var(--text);font-size:0.65rem;border-radius:4px;padding:1px 4px;cursor:pointer;outline:none}
.line-count{margin-left:auto;color:var(--muted2)}

/* Players */
.player-card{background:var(--card2);border:1px solid var(--border2);border-radius:6px;padding:10px 12px;display:flex;align-items:center;gap:10px}
.player-avatar{width:30px;height:30px;border-radius:5px;background:var(--card);border:1px solid var(--border2);display:flex;align-items:center;justify-content:center;font-size:0.7rem;color:var(--muted);font-family:var(--mono);flex-shrink:0}
.player-name{font-size:0.78rem;font-family:var(--mono)}
.player-status{font-size:0.6rem;color:var(--green);margin-top:1px}

/* File browser */
.breadcrumb{display:flex;align-items:center;gap:4px;font-family:var(--mono);font-size:0.72rem;color:var(--muted);flex-shrink:0;margin-bottom:8px;flex-wrap:wrap}
.breadcrumb-part{color:var(--text);cursor:pointer;padding:2px 4px;border-radius:3px;transition:background .12s}
.breadcrumb-part:hover{background:var(--card2)}
.breadcrumb-sep{color:var(--muted2)}
.file-list{flex:1;overflow-y:auto;border:1px solid var(--border);border-radius:6px;background:#0d0d0d}
.file-list::-webkit-scrollbar{width:4px}
.file-list::-webkit-scrollbar-thumb{background:var(--border2);border-radius:2px}
.file-item{display:flex;align-items:center;gap:10px;padding:8px 12px;border-bottom:1px solid var(--border);cursor:pointer;transition:background .12s;font-size:0.78rem}
.file-item:last-child{border-bottom:none}
.file-item:hover{background:var(--card2)}
.file-icon{font-size:14px;flex-shrink:0;width:18px;text-align:center}
.file-name{font-family:var(--mono);color:var(--text);flex:1}
.file-name.dir{color:var(--blue)}
.file-empty{padding:32px;text-align:center;color:var(--muted2);font-size:0.78rem}

/* File viewer */
.file-viewer{flex:1;display:flex;flex-direction:column;overflow:hidden}
.file-viewer-toolbar{display:flex;align-items:center;gap:8px;margin-bottom:8px;flex-shrink:0}
.file-viewer-name{font-family:var(--mono);font-size:0.78rem;color:var(--text);flex:1}
.file-textarea{flex:1;background:#0d0d0d;border:1px solid var(--border);border-radius:6px;padding:12px;font-family:var(--mono);font-size:0.72rem;color:#ccc;outline:none;resize:none;line-height:1.75}
.file-textarea:focus{border-color:#444}

/* Config */
.config-toolbar{display:flex;align-items:center;gap:8px;margin-bottom:8px;flex-shrink:0}
.config-filename{font-family:var(--mono);font-size:0.75rem;color:var(--muted)}
.config-textarea{flex:1;background:#0d0d0d;border:1px solid var(--border);border-radius:6px;padding:12px;font-family:var(--mono);font-size:0.72rem;color:#ccc;outline:none;resize:none;line-height:1.75}
.config-textarea:focus{border-color:#444}
.save-btn{background:var(--card2);border:1px solid var(--border2);color:var(--text);padding:5px 14px;border-radius:4px;font-size:0.75rem;font-weight:500;cursor:pointer;transition:border-color .12s,color .12s}
.save-btn:hover{border-color:var(--green);color:var(--green)}

/* Mods */
.mod-header{display:flex;justify-content:space-between;align-items:center;flex-shrink:0;margin-bottom:8px}
.mod-count{font-size:0.72rem;color:var(--muted)}
.mod-grid{flex:1;overflow-y:auto;display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:8px;align-content:start}
.mod-grid::-webkit-scrollbar{width:4px}
.mod-grid::-webkit-scrollbar-thumb{background:var(--border2);border-radius:2px}
.mod-item{background:var(--card2);border:1px solid var(--border2);border-radius:6px;padding:8px 10px;display:flex;flex-direction:column;gap:6px;position:relative}
.mod-icon{font-size:16px;flex-shrink:0}
.mod-name{font-family:var(--mono);font-size:0.72rem;color:var(--text);overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.mod-actions{display:flex;flex-direction:row;gap:4px;width:auto}
.mod-item:hover .mod-actions{display:flex}
.mod-btn{background:var(--card);border:1px solid var(--border2);color:var(--muted);padding:3px 6px;border-radius:3px;cursor:pointer;font-size:0.6rem;transition:color .12s,border-color .12s;text-align:center}
.mod-btn:hover{color:var(--text);border-color:#444}
.mod-btn-remove:hover{color:var(--red);border-color:var(--red)}
.mod-btn-disable:hover{color:var(--amber);border-color:var(--amber)}

/* Overview */
.overview-grid{display:grid;grid-template-columns:1fr 1fr;gap:10px;overflow-y:auto}
.ov-card{background:var(--card);border:1px solid var(--border);border-radius:8px;padding:14px}
.ov-card-title{font-size:0.65rem;font-weight:500;color:var(--muted);letter-spacing:.08em;text-transform:uppercase;margin-bottom:10px}
.ov-stat{display:flex;justify-content:space-between;padding:5px 0;border-bottom:1px solid var(--border);font-size:0.72rem}
.ov-stat:last-child{border-bottom:none}
.ov-stat-k{color:var(--muted)}
.ov-stat-v{font-family:var(--mono);font-size:0.68rem}

/* Toast */
#toast{position:fixed;bottom:18px;right:18px;z-index:9999;background:var(--card2);border:1px solid var(--border2);border-radius:6px;padding:8px 14px;font-size:0.75rem;color:var(--text);transform:translateY(50px);opacity:0;transition:all .2s ease;max-width:280px}
#toast.show{transform:translateY(0);opacity:1}
#toast.err{border-color:var(--red);color:var(--red)}
#toast.ok{border-color:var(--green);color:var(--green)}

@keyframes fadeUp{from{opacity:0;transform:translateY(6px)}to{opacity:1;transform:none}}
.stat-card{animation:fadeUp .3s ease both}
.stat-card:nth-child(1){animation-delay:.00s}
.stat-card:nth-child(2){animation-delay:.04s}
.stat-card:nth-child(3){animation-delay:.08s}
.stat-card:nth-child(4){animation-delay:.12s}

/* Mobile responsive - hide sidebar, make cards smaller */
@media (max-width: 768px) {
  .sidebar {
    display: none;
  }
  
  .app {
    flex-direction: column;
  }
  
  .stats-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 8px;
    margin-bottom: 10px;
    max-height: 250px;
    overflow-y: auto;
    padding-right: 5px;
  }
  
  .stat-card {
    padding: 8px;
    min-height: auto;
  }
  
  .stat-card-header {
    margin-bottom: 6px;
  }
  
  .stat-card-label {
    font-size: 0.55rem;
  }
  
  .stat-card-icon {
    width: 20px;
    height: 20px;
    font-size: 10px;
  }
  
  .waveform {
    height: 24px;
    margin-bottom: 4px;
  }
  
  .stat-status-text {
    font-size: 0.85rem;
  }
  
  .perf-item {
    margin-bottom: 4px;
  }
  
  .info-row {
    padding: 2px 0;
  }
  
  .content-wrap {
    flex: 1;
    overflow: hidden;
  }
  
.tab-panel {
  display: none;
  flex: 1;
  flex-direction: column;
  overflow: hidden;
  padding: 14px 22px 10px;
}

#panel-settings .settings-content {
  padding: 14px 22px;
}
  
  .console-input-row {
    margin-bottom: 10px;
  }
  
  #console {
    font-size: 0.65rem;
    padding: 8px;
  }
}

/* Desktop - show sidebar normally */
@media (min-width: 769px) {
  .sidebar {
    display: flex;
  }
}

.settings-tab-btn {
  transition: color .2s, background .2s;
}
.settings-tab-btn:hover {
  background: rgba(255,255,255,.02);
}
.settings-tab-btn.active {
  color: var(--text);
  border-right-color: var(--text);
}

.settings-content {
  display: none;
}
.settings-content.active {
  display: block;
}
.prop-toggle {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  padding: 12px;
  background: var(--card2);
  border: 1px solid var(--border2);
  border-radius: 6px;
}

.prop-toggle input[type="checkbox"] {
  cursor: pointer;
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  width: 20px;
  height: 20px;
  border: 2px solid var(--border2);
  border-radius: 6px;
  background: var(--card);
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.prop-toggle input[type="checkbox"]:hover {
  border-color: var(--green);
}

.prop-toggle input[type="checkbox"]:checked {
  background: var(--green);
  border-color: var(--green);
  background-image: url("data:image/svg+xml,%3csvg viewBox='0 0 16 16' fill='%23166534' xmlns='http://www.w3.org/2000/svg'%3e%3cpath d='M12.207 4.793a1 1 0 010 1.414l-5 5a1 1 0 01-1.414 0l-2-2a1 1 0 011.414-1.414L6.5 9.086l4.293-4.293a1 1 0 011.414 0z'/%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: center;
  background-size: 100%;
}

.prop-toggle input[type="checkbox"]:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.2);
}

.prop-key {
  font-family: var(--mono);
  font-size: 0.75rem;
  color: var(--muted);
  width: 100%;
}

.prop-value {
  width: 100%;
  padding: 6px 8px;
  background: var(--card);
  border: 1px solid var(--border2);
  border-radius: 4px;
  font-family: var(--mono);
  font-size: 0.7rem;
  color: var(--text);
}

@media (max-width: 768px) {
  #properties-list {
    grid-template-columns: repeat(2, 1fr) !important;
  }
}
</style>
</head>
<body>

<div id="login-overlay">
  <div class="login-box">
    <div class="login-logo">
      <div class="login-logo-icon"><img src="https://cravatar.eu/helmavatar/miyoubi/26" style="width:26px;height:26px;border-radius:4px"></div>
      <div class="login-logo-text">Miyoubi<span> Panel</span></div>
    </div>
    <div class="login-title">Sign in</div>
    <div class="login-sub">Enter your credentials to access the panel</div>
    <div class="form-group">
      <label class="form-label">Username</label>
      <input class="form-input" id="login-user" type="text" placeholder="admin" autocomplete="username">
    </div>
    <div class="form-group">
      <label class="form-label">Password</label>
      <input class="form-input" id="login-pass" type="password" placeholder="&#x2022;&#x2022;&#x2022;&#x2022;&#x2022;&#x2022;&#x2022;&#x2022;" autocomplete="current-password"
             onkeydown="if(event.key==='Enter') doLogin()">
    </div>
    <button class="login-btn" onclick="doLogin()">Sign in</button>
    <div class="login-error" id="login-error">Invalid username or password</div>
  </div>
</div>

<div class="app" id="app" style="display:none">
  <aside class="sidebar">
    <div class="sb-logo">
      <div class="sb-logo-icon"><img src="https://cravatar.eu/helmavatar/miyoubi/26" style="width:26px;height:26px;border-radius:4px"></div>
      <div class="sb-logo-text">Miyoubi<span> Panel</span></div>
    </div>
    <div class="sb-section">Navigation</div>
    <a class="nav-item active" id="nav-dashboard" onclick="goTab('console');setNav(this)">
      <svg class="nav-icon" viewBox="0 0 16 16" fill="currentColor"><path d="M6 9a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7A.5.5 0 0 1 6 9m-1.146-2.854a.5.5 0 0 1 0 .708L3.707 8l1.147 1.146a.5.5 0 0 1-.708.708l-1.5-1.5a.5.5 0 0 1 0-.708l1.5-1.5a.5.5 0 0 1 .708 0z"/><path d="M1 13.5A1.5 1.5 0 0 0 2.5 15h11a1.5 1.5 0 0 0 1.5-1.5v-11A1.5 1.5 0 0 0 13.5 1h-11A1.5 1.5 0 0 0 1 2.5zm1.5-12h11a.5.5 0 0 1 .5.5v11a.5.5 0 0 1-.5.5h-11a.5.5 0 0 1-.5-.5v-11a.5.5 0 0 1 .5-.5"/></svg>
      Dashboard
    </a>
    <a class="nav-item" id="nav-overview" onclick="goTab('overview');setNav(this)">
      <svg class="nav-icon" viewBox="0 0 16 16" fill="currentColor"><path d="M1 2.5A1.5 1.5 0 0 1 2.5 1h3A1.5 1.5 0 0 1 7 2.5v3A1.5 1.5 0 0 1 5.5 7h-3A1.5 1.5 0 0 1 1 5.5zM9 2.5A1.5 1.5 0 0 1 10.5 1h3A1.5 1.5 0 0 1 15 2.5v3A1.5 1.5 0 0 1 13.5 7h-3A1.5 1.5 0 0 1 9 5.5zm0 6.5A1.5 1.5 0 0 1 10.5 7.5h3A1.5 1.5 0 0 1 15 9v3a1.5 1.5 0 0 1-1.5 1.5h-3A1.5 1.5 0 0 1 9 12zM1 9a1.5 1.5 0 0 1 1.5-1.5h3A1.5 1.5 0 0 1 7 9v3a1.5 1.5 0 0 1-1.5 1.5h-3A1.5 1.5 0 0 1 1 12z"/></svg>
      Overview
    </a>
    <div class="sb-section" style="margin-top:6px">Servers</div>
    <div class="qa-item active">
      <div class="qa-dot" id="qa-dot"></div>
      <div class="qa-name" id="qa-name">My Server</div>
    </div>
    <div class="sb-spacer"></div>
    <div class="sb-user" onclick="doLogout()">
      <div class="sb-avatar" id="sb-avatar">A</div>
      <div>
        <div class="sb-user-name" id="sb-uname">Admin</div>
        <div class="sb-user-role">Panel Admin</div>
      </div>
      <div class="sb-logout">Sign out</div>
    </div>
    <div class="sb-version">v0.0.1</div>
  </aside>

  <div class="main">
    <div class="srv-header">
     <div class="srv-icon"><img src="https://cravatar.eu/helmavatar/miyoubi/44" style="width:44px;height:44px;border-radius:8px"></div>
      <div>
        <div class="srv-title" id="srv-name">Loading...</div>
        <div class="srv-desc">Docker Minecraft Server</div>
        <div class="srv-meta" id="srv-meta">Fetching status...</div>
      </div>
      <div class="srv-actions">
        <button class="act-btn act-btn-start" id="btn-start" onclick="action('start')" disabled>
          <svg width="10" height="10" viewBox="0 0 10 10" fill="currentColor"><polygon points="2,1 9,5 2,9"/></svg>Start
        </button>
        <button class="act-btn act-btn-stop" id="btn-stop" onclick="action('stop')" disabled>
          <svg width="10" height="10" viewBox="0 0 10 10" fill="currentColor"><rect x="2" y="2" width="6" height="6"/></svg>Stop
        </button>
        <button class="act-btn act-btn-restart" id="btn-restart" onclick="action('restart')" disabled>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/></svg>Restart
        </button>
      </div>
    </div>

    <div class="content-wrap">
      <div class="content-inner">
        <div class="stats-row">
          <div class="stat-card">
            <div class="stat-card-header">
              <div><div class="stat-card-label">Server Status</div><div class="stat-card-sub">Live monitoring</div></div>
              <div class="stat-card-icon">&#x2665;</div>
            </div>
            <div class="waveform" id="waveform">
              <div class="w-bar"></div><div class="w-bar"></div><div class="w-bar"></div>
              <div class="w-bar"></div><div class="w-bar"></div><div class="w-bar"></div>
            </div>
            <div class="stat-status-text" id="status-text">LOADING</div>
            <div class="stat-status-sub" id="status-sub">Connecting...</div>
          </div>
          <div class="stat-card">
            <div class="stat-card-header">
              <div><div class="stat-card-label">Connection</div><div class="stat-card-sub">Server address</div></div>
              <div class="stat-card-icon">&#x1F310;</div>
            </div>
            <div class="conn-address" onclick="copyAddr()">
              <span id="conn-addr" style="font-family:var(--mono)">localhost:25565</span>
              <span class="copy-label">Copy</span>
            </div>
            <div class="conn-hint" id="uptime-display">&nbsp;</div>
          </div>
          <div class="stat-card">
            <div class="stat-card-header">
              <div><div class="stat-card-label">Server Info</div><div class="stat-card-sub">Details &amp; versions</div></div>
              <div class="stat-card-icon">&#x2139;</div>
            </div>
            <div class="info-row"><span class="info-key">Container</span><span class="info-val" id="info-cid">-</span></div>
            <div class="info-row"><span class="info-key">Image</span><span class="info-val" id="info-image">-</span></div>
            <div class="info-row"><span class="info-key">Status</span><span class="info-val" id="info-status">-</span></div>
            <div class="info-row"><span class="info-key">Uptime</span><span class="info-val" id="info-uptime">-</span></div>
          </div>
          <div class="stat-card">
            <div class="stat-card-header">
              <div><div class="stat-card-label">Performance</div><div class="stat-card-sub">Resources &amp; metrics</div></div>
              <div class="stat-card-icon">&#x25A6;</div>
            </div>
            <div class="perf-item">
              <div class="perf-header"><span class="perf-label">Memory</span><span class="perf-val" id="perf-mem-val">-</span></div>
              <div class="perf-bar-track"><div class="perf-bar-fill" id="perf-mem-bar"></div></div>
            </div>
            <div class="perf-item">
              <div class="perf-header"><span class="perf-label">CPU</span><span class="perf-val" id="perf-cpu-val">-</span></div>
              <div class="perf-bar-track"><div class="perf-bar-fill" id="perf-cpu-bar"></div></div>
            </div>
            <div class="perf-item" style="margin-top:4px">
              <div class="perf-header"><span class="perf-label">Players</span><span class="perf-val" id="perf-players">-</span></div>
              <div id="player-pips" style="display:flex;flex-wrap:wrap;gap:4px;margin-top:4px"></div>
            </div>
          </div>
        </div>
      </div>

      <div class="tabs-bar">
        <button class="tab-btn" id="tab-overview" onclick="goTab('overview')">Overview</button>
        <button class="tab-btn active" id="tab-console" onclick="goTab('console')">Console</button>
        <button class="tab-btn" id="tab-players" onclick="goTab('players')">Players</button>
        <button class="tab-btn" id="tab-files" onclick="goTab('files')">Files</button>
        <button class="tab-btn" id="tab-mods" onclick="goTab('mods')">Mods</button>
        <button class="tab-btn" id="tab-config" onclick="goTab('config')">Config</button>
	<button class="tab-btn" id="tab-settings" onclick="goTab('settings')">Settings</button>
      </div>

      <div class="tab-content">

        <!-- Overview -->
        <div class="tab-panel" id="panel-overview">
          <div class="overview-grid">
            <div class="ov-card">
              <div class="ov-card-title">Container Details</div>
              <div class="ov-stat"><span class="ov-stat-k">Container ID</span><span class="ov-stat-v" id="ov-cid">-</span></div>
              <div class="ov-stat"><span class="ov-stat-k">Image</span><span class="ov-stat-v" id="ov-image">-</span></div>
              <div class="ov-stat"><span class="ov-stat-k">Status</span><span class="ov-stat-v" id="ov-status">-</span></div>
              <div class="ov-stat"><span class="ov-stat-k">Uptime</span><span class="ov-stat-v" id="ov-uptime">-</span></div>
            </div>
            <div class="ov-card">
              <div class="ov-card-title">Resource Usage</div>
              <div class="ov-stat"><span class="ov-stat-k">CPU</span><span class="ov-stat-v" id="ov-cpu">-</span></div>
              <div class="ov-stat"><span class="ov-stat-k">Memory Used</span><span class="ov-stat-v" id="ov-memused">-</span></div>
              <div class="ov-stat"><span class="ov-stat-k">Memory Limit</span><span class="ov-stat-v" id="ov-memlim">-</span></div>
              <div class="ov-stat"><span class="ov-stat-k">Memory %</span><span class="ov-stat-v" id="ov-mempct">-</span></div>
            </div>
          </div>
        </div>

        <!-- Console -->
        <div class="tab-panel active" id="panel-console">
          <div class="console-toolbar">
            <div class="console-title">
              <svg width="13" height="13" viewBox="0 0 16 16" fill="currentColor" style="color:#666"><path d="M6 9a.5.5 0 0 1 .5-.5h7a.5.5 0 0 1 0 1h-7A.5.5 0 0 1 6 9m-1.146-2.854a.5.5 0 0 1 0 .708L3.707 8l1.147 1.146a.5.5 0 0 1-.708.708l-1.5-1.5a.5.5 0 0 1 0-.708l1.5-1.5a.5.5 0 0 1 .708 0z"/><path d="M1 13.5A1.5 1.5 0 0 0 2.5 15h11a1.5 1.5 0 0 0 1.5-1.5v-11A1.5 1.5 0 0 0 13.5 1h-11A1.5 1.5 0 0 0 1 2.5zm1.5-12h11a.5.5 0 0 1 .5.5v11a.5.5 0 0 1-.5.5h-11a.5.5 0 0 1-.5-.5v-11a.5.5 0 0 1 .5-.5"/></svg>
              Server Console <span class="status-badge" id="console-badge">offline</span>
            </div>
            <div class="console-actions">
              
              <button class="con-btn" onclick="clearConsole()">Clear</button>
              <button class="con-btn" id="scroll-btn" onclick="toggleAutoScroll()">Auto-scroll: ON</button>
            </div>
          </div>
          <div id="console"></div>
          <div class="console-input-row">
            <div class="input-prompt">$</div>
            <input id="cmd-input" type="text" placeholder="Enter command..." autocomplete="off"
                   onkeydown="if(event.key==='Enter') sendCommand()">
            <button class="cmd-send" onclick="sendCommand()">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/></svg>
            </button>
          </div>
          <div class="console-offline-msg" id="console-offline-msg"><span>&#x25CF;</span> Server is offline &mdash; start the server to send commands</div><div class="console-footer">
            <label><input type="checkbox" id="autoscroll-cb" checked onchange="autoScroll=this.checked"> Auto-scroll</label>
            <span>Tail: <select id="tail-sel"><option>100</option><option selected>500</option><option>1000</option></select></span>
            <span class="line-count" id="line-count">0 lines</span>
          </div>
        </div>

<!-- Settings -->
<div class="tab-panel" id="panel-settings">
  <div style="display:flex;height:100%;width:100%;gap:0;margin:0">
    <!-- Settings sidebar -->
    <div style="width:160px;border-right:1px solid var(--border);overflow-y:auto;flex-shrink:0;background:var(--card)">
      <button class="settings-tab-btn active" onclick="goSettingsTab('server')" style="width:100%;text-align:left;padding:12px;border:none;background:transparent;color:var(--text);cursor:pointer;border-right:2px solid transparent;border-right-color:var(--text);font-size:0.8rem;font-weight:500;margin:0">Server</button>
      <button class="settings-tab-btn" onclick="goSettingsTab('resources')" style="width:100%;text-align:left;padding:12px;border:none;background:transparent;color:var(--muted);cursor:pointer;border-right:2px solid transparent;font-size:0.8rem;font-weight:500;margin:0">Resources</button>
      <button class="settings-tab-btn" onclick="goSettingsTab('properties')" style="width:100%;text-align:left;padding:12px;border:none;background:transparent;color:var(--muted);cursor:pointer;border-right:2px solid transparent;font-size:0.8rem;font-weight:500;margin:0">Properties</button>
    </div>
    
    <!-- Settings content scroll wrapper -->
    <div style="flex:1;overflow-y:auto;width:100%">
      <!-- Server Tab -->
      <div id="settings-server" class="settings-content active" style="padding:14px 22px">
        <div style="margin-bottom:20px">
          <h3 style="font-size:0.9rem;margin-bottom:10px;color:var(--text)">Server Options</h3>
          <div style="display:grid;grid-template-columns:1fr 1fr 1fr;gap:10px">
            <div>
              <label style="font-size:0.75rem;color:var(--muted);display:block;margin-bottom:6px">Difficulty</label>
              <select id="setting-difficulty" style="width:100%;padding:6px;background:var(--card2);border:1px solid var(--border2);color:var(--text);border-radius:4px">
                <option value="peaceful">Peaceful</option>
                <option value="easy">Easy</option>
                <option value="normal">Normal</option>
                <option value="hard">Hard</option>
              </select>
            </div>
            <div>
              <label style="font-size:0.75rem;color:var(--muted);display:block;margin-bottom:6px">Game Mode</label>
              <select id="setting-mode" style="width:100%;padding:6px;background:var(--card2);border:1px solid var(--border2);color:var(--text);border-radius:4px">
                <option value="survival">Survival</option>
                <option value="creative">Creative</option>
                <option value="adventure">Adventure</option>
                <option value="spectator">Spectator</option>
              </select>
            </div>
            <div style="display:flex;align-items:flex-end">
              <button class="save-btn" onclick="applyServerSettings()" style="width:100%">Apply</button>
            </div>
          </div>
        </div>

        <div style="margin-bottom:20px">
          <h3 style="font-size:0.9rem;margin-bottom:10px;color:var(--text)">Performance</h3>
          <div style="display:flex;align-items:flex-end;gap:10px">
            <div style="flex:1">
              <label style="font-size:0.75rem;color:var(--muted);display:block;margin-bottom:6px">Memory (e.g. 2G, 4G)</label>
              <input id="setting-memory" type="text" placeholder="1G" style="width:100%;padding:6px;background:var(--card2);border:1px solid var(--border2);color:var(--text);border-radius:4px;font-family:var(--mono)">
            </div>
            <button class="save-btn" onclick="applyServerSettings()" style="align-self:flex-end">Apply</button>
          </div>
        </div>

        <div style="margin-bottom:20px">
          <h3 style="font-size:0.9rem;margin-bottom:10px;color:var(--text)">Features</h3>
          <div style="display:flex;flex-direction:column;gap:10px">
            <label style="display:flex;align-items:center;gap:8px;cursor:pointer">
              <input type="checkbox" id="setting-autopause" style="cursor:pointer">
              <span style="font-size:0.8rem">Enable AutoPause</span>
            </label>
            <label style="display:flex;align-items:center;gap:8px;cursor:pointer">
              <input type="checkbox" id="setting-autostop" style="cursor:pointer">
              <span style="font-size:0.8rem">Enable AutoStop</span>
            </label>
            <label style="display:flex;align-items:center;gap:8px;cursor:pointer">
              <input type="checkbox" id="setting-whitelist" style="cursor:pointer">
              <span style="font-size:0.8rem">Enable Whitelist</span>
            </label>
          </div>
          <button class="save-btn" onclick="saveFeatures()" style="margin-top:10px;width:100%">Save Features</button>
        </div>

        <div style="margin-bottom:20px">
          <h3 style="font-size:0.9rem;margin-bottom:10px;color:var(--text)">Server Message (MOTD)</h3>
          <textarea id="setting-motd" placeholder="A Minecraft Server" style="width:100%;padding:8px;background:var(--card2);border:1px solid var(--border2);color:var(--text);border-radius:4px;font-family:var(--font);min-height:80px;resize:none;font-size:0.75rem"></textarea>
          <button class="save-btn" onclick="saveSetting('MOTD')" style="margin-top:8px;width:100%">Update MOTD</button>
        </div>
      </div>

<!-- Resources Tab -->
<div id="settings-resources" class="settings-content" style="display:none;padding:14px 22px">
  <h3 style="font-size:0.9rem;margin-bottom:15px;color:var(--text)">Docker Resources</h3>
  <div id="resources-list" style="display:grid;grid-template-columns:repeat(auto-fill,minmax(240px,1fr));gap:12px;align-content:start">
    <div class="prop-toggle"><input type="checkbox" class="res-checkbox" data-key="cpu-limit"><span class="prop-key">CPU Limit</span><input type="text" class="prop-value" data-key="cpu-limit" value="2" placeholder="2"></div>
    <div class="prop-toggle"><input type="checkbox" class="res-checkbox" data-key="memory-limit"><span class="prop-key">Memory Limit</span><input type="text" class="prop-value" data-key="memory-limit" value="10G" placeholder="10G"></div>
    <div class="prop-toggle"><input type="checkbox" class="res-checkbox" data-key="cpu-reservation"><span class="prop-key">CPU Reservation</span><input type="text" class="prop-value" data-key="cpu-reservation" value="0.3" placeholder="0.3"></div>
    <div class="prop-toggle"><input type="checkbox" class="res-checkbox" data-key="memory-reservation"><span class="prop-key">Memory Reservation</span><input type="text" class="prop-value" data-key="memory-reservation" value="4G" placeholder="4G"></div>
  </div>
  <button class="save-btn" onclick="applyResources()" style="width:100%;margin-top:20px">Apply Resources</button>
</div>        

<!-- Properties Tab -->
<div id="settings-properties" class="settings-content" style="display:none;padding:14px 22px">
  <h3 style="font-size:0.9rem;margin-bottom:15px;color:var(--text)">Server Properties (Override)</h3>
  <div id="properties-list" style="display:grid;grid-template-columns:repeat(auto-fill,minmax(240px,1fr));gap:12px;align-content:start">
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="difficulty"><span class="prop-key">difficulty</span><input type="text" class="prop-value" data-key="difficulty" value="hard"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="gamemode"><span class="prop-key">gamemode</span><input type="text" class="prop-value" data-key="gamemode" value="survival"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="level-name"><span class="prop-key">level-name</span><input type="text" class="prop-value" data-key="level-name" value="world"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="max-players"><span class="prop-key">max-players</span><input type="text" class="prop-value" data-key="max-players" value="20"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="motd"><span class="prop-key">motd</span><input type="text" class="prop-value" data-key="motd" value="A Minecraft Server"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="online-mode"><span class="prop-key">online-mode</span><input type="text" class="prop-value" data-key="online-mode" value="false"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="pvp"><span class="prop-key">pvp</span><input type="text" class="prop-value" data-key="pvp" value="true"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="spawn-protection"><span class="prop-key">spawn-protection</span><input type="text" class="prop-value" data-key="spawn-protection" value="16"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="view-distance"><span class="prop-key">view-distance</span><input type="text" class="prop-value" data-key="view-distance" value="10"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="simulation-distance"><span class="prop-key">simulation-distance</span><input type="text" class="prop-value" data-key="simulation-distance" value="10"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="white-list"><span class="prop-key">white-list</span><input type="text" class="prop-value" data-key="white-list" value="false"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="enable-rcon"><span class="prop-key">enable-rcon</span><input type="text" class="prop-value" data-key="enable-rcon" value="false"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="rcon.port"><span class="prop-key">rcon.port</span><input type="text" class="prop-value" data-key="rcon.port" value="25575"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="allow-flight"><span class="prop-key">allow-flight</span><input type="text" class="prop-value" data-key="allow-flight" value="false"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="enforce-whitelist"><span class="prop-key">enforce-whitelist</span><input type="text" class="prop-value" data-key="enforce-whitelist" value="false"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="level-seed"><span class="prop-key">level-seed</span><input type="text" class="prop-value" data-key="level-seed" value=""></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="enable-query"><span class="prop-key">enable-query</span><input type="text" class="prop-value" data-key="enable-query" value="false"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="query.port"><span class="prop-key">query.port</span><input type="text" class="prop-value" data-key="query.port" value="22553"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="server-port"><span class="prop-key">server-port</span><input type="text" class="prop-value" data-key="server-port" value="25565"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="server-ip"><span class="prop-key">server-ip</span><input type="text" class="prop-value" data-key="server-ip" value=""></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="max-tick-time"><span class="prop-key">max-tick-time</span><input type="text" class="prop-value" data-key="max-tick-time" value="60000"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="op-permission-level"><span class="prop-key">op-permission-level</span><input type="text" class="prop-value" data-key="op-permission-level" value="4"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="function-permission-level"><span class="prop-key">function-permission-level</span><input type="text" class="prop-value" data-key="function-permission-level" value="2"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="entity-broadcast-range-percentage"><span class="prop-key">entity-broadcast-range-percentage</span><input type="text" class="prop-value" data-key="entity-broadcast-range-percentage" value="100"></div>
    <div class="prop-toggle"><input type="checkbox" class="prop-checkbox" data-key="rate-limit"><span class="prop-key">rate-limit</span><input type="text" class="prop-value" data-key="rate-limit" value="0"></div>
  </div>
  <button class="save-btn" onclick="applyProperties()" style="width:100%;margin-top:20px">Apply Properties</button>
</div>
</div>  <!-- This closes the settings-content scroll wrapper -->
</div>  <!-- This closes the settings flex container -->
</div>  <!-- This closes panel-settings -->

        <!-- Players -->
        <div class="tab-panel" id="panel-players">
          <div style="display:flex;align-items:center;justify-content:space-between;flex-shrink:0;margin-bottom:10px">
            <div style="font-size:0.72rem;color:var(--muted)">Online players <span id="pl-count-label" style="color:var(--text);font-family:var(--mono)">-</span></div>
            <button class="con-btn" onclick="fetchPlayers()">Refresh</button>
          </div>
          <div style="flex:1;overflow-y:auto">
            <div id="player-list-empty" style="text-align:center;padding:48px 0;color:var(--muted2);font-size:0.78rem">No players online</div>
            <div id="player-list" style="display:grid;grid-template-columns:repeat(auto-fill,minmax(180px,1fr));gap:8px"></div>
          </div>
        </div>

        <!-- Files -->
        <div class="tab-panel" id="panel-files">
          <div class="breadcrumb" id="file-breadcrumb"></div>
          <div id="file-list-wrap" style="flex:1;display:flex;flex-direction:column;overflow:hidden">
            <div class="file-list" id="file-list"></div>
          </div>
          <div id="file-viewer-wrap" class="file-viewer" style="display:none">
            <div class="file-viewer-toolbar">
              <button class="con-btn" onclick="closeFileViewer()">&#x2190; Back</button>
              <span class="file-viewer-name" id="viewer-filename"></span>
              <button class="save-btn" id="viewer-save-btn" onclick="saveFile()">Save</button>
            </div>
            <textarea class="file-textarea" id="viewer-content" spellcheck="false"></textarea>
          </div>
        </div>

	<!-- Mods -->
	  <div class="tab-panel" id="panel-mods">
  	    <div class="mod-header">
    	    <div class="mod-count" id="mod-count">Loading mods...</div>
    	    <div style="margin-left:auto;display:flex;gap:6px">
      	    <button class="con-btn" id="disabled-toggle-btn" onclick="toggleDisabledView()" style="display:none">Show Disabled</button>
      	    <button class="con-btn" onclick="fetchMods()">Refresh</button>
    	</div>
  	</div>
 	 <div class="mod-grid" id="mod-grid"></div>
	</div>

        <!-- Config -->
        <div class="tab-panel" id="panel-config">
          <div class="config-toolbar">
            <span class="config-filename" id="config-filename">docker-compose.yml</span>
            <div style="margin-left:auto;display:flex;gap:6px;align-items:center">
              <span id="config-status" style="font-size:0.65rem;color:var(--muted)"></span>
              <button class="save-btn" onclick="saveConfig()">Save</button>
            </div>
          </div>
          <textarea class="config-textarea" id="config-content" spellcheck="false" placeholder="Loading..."></textarea>
        </div>

      </div>
    </div>
  </div>
</div>

<div id="toast"></div>

<script>
var API = '';
var autoScroll = true;
var lineCount = 0;
var sseSource = null;
var currentUser = null;
var pollTimer, playerPollTimer;

// Credentials — change these or extend with /api/login
var USERS = { 'admin': 'admin' };

// Cookie helpers
function setCookie(name, val, days) {
  var exp = new Date(Date.now() + days * 864e5).toUTCString();
  document.cookie = name + '=' + encodeURIComponent(val) + '; expires=' + exp + '; path=/; SameSite=Strict';
}
function getCookie(name) {
  var m = document.cookie.match('(?:^|;)\s*' + name + '=([^;]*)');
  return m ? decodeURIComponent(m[1]) : null;
}
function clearCookie(name) {
  document.cookie = name + '=; max-age=0; path=/';
}

// Auto-login from cookie on page load
window.addEventListener('load', function() {
  var saved = getCookie('mcpanel_user');
  if (saved && USERS[saved]) {
    currentUser = saved;
    showApp(saved);
  }
});

function doLogin() {
  var u = document.getElementById('login-user').value.trim();
  var p = document.getElementById('login-pass').value;
  var err = document.getElementById('login-error');
  if (USERS[u] && USERS[u] === p) {
    currentUser = u;
    setCookie('mcpanel_user', u, 7);
    showApp(u);
    err.style.display = 'none';
  } else {
    err.style.display = 'block';
    document.getElementById('login-pass').value = '';
    document.getElementById('login-pass').focus();
  }
}

function showApp(u) {
  document.getElementById('login-overlay').style.display = 'none';
  document.getElementById('app').style.display = 'flex';
  document.getElementById('sb-avatar').textContent = u[0].toUpperCase();
  document.getElementById('sb-uname').textContent = u;
  fetchStatus();
  loadHistoricalLogs();
  startLogStream();
  fetchPlayers();
}

function loadHistoricalLogs() {
  var logs = JSON.parse(localStorage.getItem('consoleLogs') || '[]');
  logs.forEach(function(log) {
    var el = document.getElementById('console');
    var div = document.createElement('div');
    div.className = log.cls;
    div.textContent = log.text;
    el.appendChild(div);
    lineCount++;
  });
  document.getElementById('line-count').textContent = lineCount + ' lines';
}

function doLogout() {
  currentUser = null;
  clearCookie('mcpanel_user');
  clearTimeout(pollTimer);
  clearTimeout(playerPollTimer);
  if (sseSource) { sseSource.close(); sseSource = null; }
  document.getElementById('app').style.display = 'none';
  document.getElementById('login-overlay').style.display = 'flex';
  document.getElementById('login-user').value = '';
  document.getElementById('login-pass').value = '';
}

// Tab switching
var tabLoaded = {};
function goTab(name) {
  document.querySelectorAll('.tab-panel').forEach(function(p) { p.classList.remove('active'); });
  document.querySelectorAll('.tab-btn').forEach(function(b) { b.classList.remove('active'); });
  var panel = document.getElementById('panel-' + name);
  var btn = document.getElementById('tab-' + name);
  if (panel) panel.classList.add('active');
  if (btn) btn.classList.add('active');
  // Lazy-load tab data
  if (!tabLoaded[name]) {
    tabLoaded[name] = true;
    if (name === 'files') loadFiles('/data');
    if (name === 'mods') fetchMods();
    if (name === 'config') fetchConfig();
  }
}
function setNav(el) {
  document.querySelectorAll('.nav-item').forEach(function(n) { n.classList.remove('active'); });
  el.classList.add('active');
}

// Status
async function fetchStatus() {
  if (!currentUser) return;
  try {
    const res = await fetch(API + '/api/status', { 
      signal: AbortSignal.timeout(5000)
    });
    if (!res.ok) throw new Error('Status API error');
    const s = await res.json();
    applyStatus(s);
    checkAndClearConsole(s);
    
    // Only keep polling if server is running
    clearTimeout(pollTimer);
    pollTimer = setTimeout(fetchStatus, 5000);
  } catch(e) { 
    setOffline();
    // Stop polling entirely on error
    clearTimeout(pollTimer);
  }
}

function applyStatus(s) {
  var r = s.running;
  document.getElementById('waveform').className = 'waveform' + (r ? ' running' : '');
  var st = document.getElementById('status-text');
  st.textContent = (s.status || 'unknown').toUpperCase();
  st.className = 'stat-status-text' + (r ? ' running' : (s.status !== 'not found' ? ' stopped' : ''));
  document.getElementById('status-sub').textContent = r ? 'Server healthy and responding' : 'Server is not running';
  var badge = document.getElementById('console-badge');
  badge.textContent = r ? 'running' : 'offline';
  badge.className = 'status-badge' + (r ? ' running' : ' stopped');
  document.getElementById('srv-name').textContent = s.container_id ? 'Container ' + s.container_id : 'Minecraft Server';
  document.getElementById('srv-meta').textContent = r ? 'Running for ' + (s.uptime || '-') : 'Server offline';
  document.getElementById('uptime-display').textContent = s.uptime ? 'Up ' + s.uptime : ' ';
  document.getElementById('qa-dot').className = 'qa-dot' + (r ? ' on' : '');
  document.getElementById('qa-name').textContent = s.container_id || 'My Server';
  document.getElementById('info-cid').textContent = s.container_id || '-';
  document.getElementById('info-image').textContent = s.image || '-';
  document.getElementById('info-status').textContent = (s.status || '-').toUpperCase();
  document.getElementById('info-uptime').textContent = s.uptime || '-';
  document.getElementById('ov-cid').textContent = s.container_id || '-';
  document.getElementById('ov-image').textContent = s.image || '-';
  document.getElementById('ov-status').textContent = (s.status || '-').toUpperCase();
  document.getElementById('ov-uptime').textContent = s.uptime || '-';
  var cpu = s.cpu_percent || 0;
  var mu = s.mem_usage_mb || 0;
  var ml = s.mem_limit_mb || 0;
  var mp = ml > 0 ? (mu / ml * 100) : 0;
  if (r) {
    document.getElementById('perf-cpu-val').textContent = cpu.toFixed(1) + '%';
    document.getElementById('perf-mem-val').textContent = mu.toFixed(0) + ' / ' + ml.toFixed(0) + ' MB';
    document.getElementById('ov-cpu').textContent = cpu.toFixed(1) + '%';
    document.getElementById('ov-memused').textContent = mu.toFixed(0) + ' MB';
    document.getElementById('ov-memlim').textContent = ml.toFixed(0) + ' MB';
    document.getElementById('ov-mempct').textContent = mp.toFixed(1) + '%';
  } else {
    ['perf-cpu-val','perf-mem-val','ov-cpu','ov-memused','ov-memlim','ov-mempct'].forEach(function(id) { document.getElementById(id).textContent = '-'; });
  }
  setBar('perf-cpu-bar', cpu);
  setBar('perf-mem-bar', mp);
  document.getElementById('btn-start').disabled = r;
  document.getElementById('btn-stop').disabled = !r;
  document.getElementById('btn-restart').disabled = !r;
  setConsoleInputEnabled(r);
}

// Auto-clear console if panel started but server isn't running
function checkAndClearConsole(s) {
  var isServerRunning = s.running;
  var isPanelJustStarted = lineCount > 0 && !document.getElementById('console').dataset.serverWasRunning;
  
  // If server is not running and this is the first check after panel start
  if (!isServerRunning && isPanelJustStarted) {
    clearConsole();
    document.getElementById('console').dataset.serverWasRunning = 'false';
  } else if (isServerRunning) {
    document.getElementById('console').dataset.serverWasRunning = 'true';
  }
}

function setConsoleInputEnabled(enabled) {
  var inp = document.getElementById('cmd-input');
  var btn = document.querySelector('.cmd-send');
  var msg = document.getElementById('console-offline-msg');
  inp.disabled = !enabled;
  inp.placeholder = enabled ? 'Enter command...' : 'Server offline';
  if (btn) btn.disabled = !enabled;
  if (msg) msg.className = 'console-offline-msg' + (enabled ? '' : ' visible');
}

function setOffline() {
  document.getElementById('status-text').textContent = 'OFFLINE';
  document.getElementById('status-text').className = 'stat-status-text stopped';
  document.getElementById('waveform').className = 'waveform';
  setConsoleInputEnabled(false);
}

function setBar(id, pct) {
  var el = document.getElementById(id);
  el.style.width = Math.min(100, pct).toFixed(1) + '%';
  el.className = 'perf-bar-fill' + (pct > 85 ? ' danger' : pct > 65 ? ' warn' : '');
}

// Actions
async function action(cmd) {
  ['btn-start','btn-stop','btn-restart'].forEach(function(id) { document.getElementById(id).disabled = true; });
  try {
    var res = await fetch(API + '/api/' + cmd, { method: 'POST' });
    var data = await res.json();
    toast(data.message || 'Done', data.ok ? 'ok' : 'err');
  } catch(e) { toast('Request failed', 'err'); }
  setTimeout(fetchStatus, 1500);
}

// Console
function startLogStream() {
  if (sseSource) sseSource.close();
  var tail = document.getElementById('tail-sel').value;
  sseSource = new EventSource(API + '/api/logs?tail=' + tail);
  sseSource.onmessage = function(e) { appendLog(e.data); };
  sseSource.onerror = function() {
    sseSource.close();
    setTimeout(function() { if (currentUser) startLogStream(); }, 5000);
  };
}

function appendLog(raw) {
  var el = document.getElementById('console');
  var div = document.createElement('div');
  var lo = raw.toLowerCase();
  var cls = 'log-line';
  if (lo.indexOf('error') !== -1 || lo.indexOf('exception') !== -1 || lo.indexOf('fatal') !== -1) cls += ' log-error';
  else if (lo.indexOf('warn') !== -1) cls += ' log-warn';
  else if (lo.indexOf('info') !== -1 || lo.indexOf('joined') !== -1 || lo.indexOf('left') !== -1) cls += ' log-info';
  div.className = cls;
  div.textContent = raw.replace(/^\d{4}-\d{2}-\d{2}T[\d:.]+Z\s/, '');
  el.appendChild(div);
  lineCount++;
  document.getElementById('line-count').textContent = lineCount + ' lines';
  
  // Save to localStorage
  var logs = JSON.parse(localStorage.getItem('consoleLogs') || '[]');
  logs.push({ cls: cls, text: div.textContent });
  logs = logs.slice(-500); // Keep last 500
  localStorage.setItem('consoleLogs', JSON.stringify(logs));
  
  while (el.children.length > 2000) { el.removeChild(el.firstChild); lineCount--; }
  if (autoScroll) el.scrollTop = el.scrollHeight;
}

function clearConsole() {
  document.getElementById('console').innerHTML = '';
  lineCount = 0;
  document.getElementById('line-count').textContent = '0 lines';
}

function toggleAutoScroll() {
  autoScroll = !autoScroll;
  document.getElementById('autoscroll-cb').checked = autoScroll;
  document.getElementById('scroll-btn').textContent = 'Auto-scroll: ' + (autoScroll ? 'ON' : 'OFF');
}

document.getElementById('tail-sel').addEventListener('change', startLogStream);

async function sendCommand() {
  var input = document.getElementById('cmd-input');
  var cmd = input.value.trim();
  if (!cmd) return;
  appendLog('> ' + cmd);
  input.value = '';
  try {
    var res = await fetch(API + '/api/command', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ command: cmd })
    });
    var data = await res.json();
    if (!data.ok) toast(data.message || 'Failed', 'err');
    else if (data.output) appendLog(data.output);
  } catch(e) { toast('Request failed', 'err'); }
}

// Players
async function fetchPlayers() {
  if (!currentUser) return;
  try {
    const res = await fetch(API + '/api/players');
    if (!res.ok) throw new Error('Players API error');
    const data = await res.json();
    applyPlayers(data);
    
    // Only keep polling if server is running
    clearTimeout(playerPollTimer);
    playerPollTimer = setTimeout(fetchPlayers, 10000);
  } catch(e) {
    // Stop polling on error
    clearTimeout(playerPollTimer);
  }
}

function applyPlayers(data) {
  var count = data.count || 0;
  var max = data.max || 20;
  var names = data.online || [];
  document.getElementById('perf-players').textContent = count + ' / ' + max;
  var pips = document.getElementById('player-pips');
  pips.innerHTML = '';
  names.forEach(function(name) {
    var pip = document.createElement('div');
    pip.title = name;
    pip.style.cssText = 'width:7px;height:7px;border-radius:50%;background:var(--green);opacity:.8';
    pips.appendChild(pip);
  });
  document.getElementById('pl-count-label').textContent = count + ' / ' + max;
  var list = document.getElementById('player-list');
  var empty = document.getElementById('player-list-empty');
  list.innerHTML = '';
  if (names.length === 0) {
    empty.style.display = 'block';
  } else {
    empty.style.display = 'none';
    names.forEach(function(name, i) {
      var card = document.createElement('div');
      card.className = 'player-card';
      card.style.animationDelay = (i * 0.04) + 's';
      card.innerHTML = '<div class="player-avatar">' + name.slice(0,2).toUpperCase() + '</div>' +
        '<div><div class="player-name">' + escHtml(name) + '</div><div class="player-status">online</div></div>';
      list.appendChild(card);
    });
  }
}

// File browser
var filePath = '/data';
var viewingFilePath = '';

function loadFiles(path) {
  filePath = path;
  document.getElementById('file-list-wrap').style.display = 'flex';
  document.getElementById('file-viewer-wrap').style.display = 'none';
  renderBreadcrumb(path);
  var list = document.getElementById('file-list');
  list.innerHTML = '<div class="file-empty">Loading...</div>';
  fetch(API + '/api/files?path=' + encodeURIComponent(path))
    .then(function(r) { return r.json(); })
    .then(function(data) { renderFileList(data); })
    .catch(function() { list.innerHTML = '<div class="file-empty">Failed to load directory</div>'; });
}

function renderBreadcrumb(path) {
  var bc = document.getElementById('file-breadcrumb');
  bc.innerHTML = '';
  var parts = path.replace(/\/+$/, '').split('/').filter(Boolean);
  var accumulated = '';
  var rootSpan = document.createElement('span');
  rootSpan.className = 'breadcrumb-part';
  rootSpan.textContent = '/';
  rootSpan.onclick = (function() { return function() { loadFiles('/data'); }; })();
  bc.appendChild(rootSpan);
  parts.forEach(function(part, i) {
    accumulated += '/' + part;
    var sep = document.createElement('span');
    sep.className = 'breadcrumb-sep';
    sep.textContent = '/';
    bc.appendChild(sep);
    var span = document.createElement('span');
    span.className = 'breadcrumb-part';
    span.textContent = part;
    var p = accumulated;
    span.onclick = function() { loadFiles(p); };
    bc.appendChild(span);
  });
}

function renderFileList(data) {
  var list = document.getElementById('file-list');
  list.innerHTML = '';
  if (!data.entries || data.entries.length === 0) {
    list.innerHTML = '<div class="file-empty">Empty directory</div>';
    return;
  }
  // Dirs first
  var sorted = data.entries.slice().sort(function(a, b) {
    if (a.is_dir === b.is_dir) return a.name.localeCompare(b.name);
    return a.is_dir ? -1 : 1;
  });
  sorted.forEach(function(entry) {
    var item = document.createElement('div');
    item.className = 'file-item';
    var icon = entry.is_dir ? '&#x1F4C1;' : getFileIcon(entry.name);
    item.innerHTML = '<span class="file-icon">' + icon + '</span>' +
      '<span class="file-name' + (entry.is_dir ? ' dir' : '') + '">' + escHtml(entry.name) + '</span>';
    if (entry.is_dir) {
      var p = filePath.replace(/\/+$/, '') + '/' + entry.name;
      item.onclick = function() { loadFiles(p); };
    } else {
      var fp = filePath.replace(/\/+$/, '') + '/' + entry.name;
      item.onclick = function() { openFile(fp, entry.name); };
    }
    list.appendChild(item);
  });
}

function getFileIcon(name) {
  var ext = name.split('.').pop().toLowerCase();
  var icons = { jar:'&#x1F4E6;', json:'&#x1F4CB;', yml:'&#x1F4CB;', yaml:'&#x1F4CB;',
    properties:'&#x2699;', txt:'&#x1F4C4;', log:'&#x1F4DC;', sh:'&#x1F4BB;', md:'&#x1F4D6;', png:'&#x1F5BC;', jpg:'&#x1F5BC;' };
  return icons[ext] || '&#x1F4C4;';
}

function openFile(path, name) {
  document.getElementById('viewer-filename').textContent = name;
  viewingFilePath = path;
  document.getElementById('viewer-content').value = 'Loading...';
  document.getElementById('file-list-wrap').style.display = 'none';
  document.getElementById('file-viewer-wrap').style.display = 'flex';
  document.getElementById('viewer-save-btn').style.display = '';
  fetch(API + '/api/files/content?path=' + encodeURIComponent(path))
    .then(function(r) { return r.json(); })
    .then(function(data) {
      if (data.binary) {
        document.getElementById('viewer-content').value = '[Binary file - cannot display]';
        document.getElementById('viewer-save-btn').style.display = 'none';
      } else {
        document.getElementById('viewer-content').value = data.content || '';
      }
    })
    .catch(function() { document.getElementById('viewer-content').value = 'Failed to load file.'; });
}

function closeFileViewer() {
  document.getElementById('file-list-wrap').style.display = 'flex';
  document.getElementById('file-viewer-wrap').style.display = 'none';
}

async function saveFile() {
  var content = document.getElementById('viewer-content').value;
  try {
    var res = await fetch(API + '/api/files/write', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ path: viewingFilePath, content: content })
    });
    var data = await res.json();
    toast(data.message || 'Saved', data.ok ? 'ok' : 'err');
  } catch(e) { toast('Save failed', 'err'); }
}

// Mods
var showingDisabled = false;

async function fetchMods() {
  document.getElementById('mod-count').textContent = 'Loading mods...';
  document.getElementById('mod-grid').innerHTML = '';
  try {
    var res = await fetch(API + '/api/mods');
    var data = await res.json();
    var entries = (data.entries || []).filter(function(e) { return !e.is_dir; });
    
    var disabledEntries = [];
    try {
      var disabledRes = await fetch(API + '/api/files?path=' + encodeURIComponent('/data/mods/disabled'));
      var disabledData = await disabledRes.json();
      disabledEntries = (disabledData.entries || []).filter(function(e) { return !e.is_dir; });
    } catch(e) {
      // disabled folder doesn't exist yet
    }
    
    document.getElementById('disabled-toggle-btn').style.display = disabledEntries.length > 0 ? 'block' : 'none';
    
    var entriesToShow = showingDisabled ? disabledEntries : entries;
    var totalCount = entries.length;
    var disabledCount = disabledEntries.length;
    
    if (showingDisabled) {
      document.getElementById('mod-count').textContent = 'Disabled mods (' + disabledCount + ')';
    } else {
      document.getElementById('mod-count').textContent = totalCount + ' mod' + (totalCount !== 1 ? 's' : '') + ' installed' + (disabledCount > 0 ? ' (' + disabledCount + ' disabled)' : '');
    }
    
    var grid = document.getElementById('mod-grid');
    if (entriesToShow.length === 0) {
      grid.innerHTML = '<div class="file-empty" style="grid-column:1/-1">No mods found</div>';
    } else {
      entriesToShow.forEach(function(e, i) {
        var item = document.createElement('div');
        item.className = 'mod-item';
        item.style.animationDelay = (i * 0.02) + 's';
        item.style.animation = 'fadeUp .25s ease both';
        var modName = escHtml(e.name);
        var modPath = (showingDisabled ? '/data/mods/disabled/' : '/data/mods/') + e.name;
        item.innerHTML = '<span class="mod-icon">&#x1F4E6;</span>' +
          '<span class="mod-name" title="' + modName + '">' + modName + '</span>' +
          '<div class="mod-actions">' +
          (showingDisabled ? 
            '<button class="mod-btn mod-btn-disable" onclick="enableMod(\'' + modPath.replace(/'/g, "\\'") + '\')">Enable</button>' +
            '<button class="mod-btn mod-btn-remove" onclick="removeMod(\'' + modPath.replace(/'/g, "\\'") + '\')">Remove</button>'
          :
            '<button class="mod-btn mod-btn-disable" onclick="disableMod(\'' + modPath.replace(/'/g, "\\'") + '\')">Disable</button>' +
            '<button class="mod-btn mod-btn-remove" onclick="removeMod(\'' + modPath.replace(/'/g, "\\'") + '\')">Remove</button>'
          ) +
          '</div>';
        grid.appendChild(item);
      });
    }
  } catch(e) {
    document.getElementById('mod-count').textContent = 'Failed to load mods';
  }
}

function toggleDisabledView() {
  showingDisabled = !showingDisabled;
  document.getElementById('disabled-toggle-btn').textContent = showingDisabled ? 'Show Active' : 'Show Disabled';
  fetchMods();
}

async function enableMod(path) {
  try {
    var res = await fetch(API + '/api/mods/enable', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ path: path })
    });
    var data = await res.json();
    toast(data.message || 'Mod enabled', data.ok ? 'ok' : 'err');
    if (data.ok) fetchMods();
  } catch(e) {
    toast('Failed to enable mod', 'err');
  }
}

async function disableMod(path) {
  try {
    var res = await fetch(API + '/api/mods/disable', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ path: path })
    });
    var data = await res.json();
    toast(data.message || 'Mod disabled', data.ok ? 'ok' : 'err');
    if (data.ok) fetchMods();
  } catch(e) {
    toast('Failed to disable mod', 'err');
  }
}

async function removeMod(path) {
  if (!confirm('Are you sure you want to remove this mod?')) return;
  try {
    var res = await fetch(API + '/api/mods/remove', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ path: path })
    });
    var data = await res.json();
    toast(data.message || 'Mod removed', data.ok ? 'ok' : 'err');
    if (data.ok) fetchMods();
  } catch(e) {
    toast('Failed to remove mod', 'err');
  }
}

// Config
async function fetchConfig() {
  document.getElementById('config-content').value = 'Loading...';
  try {
    var res = await fetch(API + '/api/config');
    var data = await res.json();
    document.getElementById('config-filename').textContent = data.filename || 'docker-compose.yml';
    document.getElementById('config-content').value = data.content || '';
  } catch(e) {
    document.getElementById('config-content').value = 'Could not load docker-compose file.\nMake sure it is in the same directory as the panel binary.';
  }
}

async function saveConfig() {
  var content = document.getElementById('config-content').value;
  var statusEl = document.getElementById('config-status');
  statusEl.textContent = 'Saving...';
  try {
    var res = await fetch(API + '/api/config', {
      method: 'POST', headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content: content })
    });
    var data = await res.json();
    statusEl.textContent = data.ok ? 'Saved' : (data.message || 'Error');
    statusEl.style.color = data.ok ? 'var(--green)' : 'var(--red)';
    setTimeout(function() { statusEl.textContent = ''; }, 3000);
    toast(data.message || 'Saved', data.ok ? 'ok' : 'err');
  } catch(e) {
    statusEl.textContent = 'Failed';
    statusEl.style.color = 'var(--red)';
  }
}

// Utils
function copyAddr() {
  navigator.clipboard.writeText(document.getElementById('conn-addr').textContent)
    .then(function() { toast('Address copied', 'ok'); });
}

function escHtml(s) {
  return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;');
}

var toastTimer;
function toast(msg, type) {
  var el = document.getElementById('toast');
  el.textContent = msg;
  el.className = 'show' + (type === 'err' ? ' err' : type === 'ok' ? ' ok' : '');
  clearTimeout(toastTimer);
  toastTimer = setTimeout(function() { el.className = ''; }, 3000);
}

// Settings functions
async function loadSettings() {
  try {
    const configRes = await fetch(API + '/api/config');
    const config = await configRes.json();
    const content = config.content || '';
    
    // Parse environment variables from docker-compose
    function getEnvValue(envVar) {
      const pattern = new RegExp(envVar + ':\\s*["\']?([^"\'\\n]+)["\']?');
      const match = content.match(pattern);
      return match ? match[1].trim() : null;
    }
    
    // Load difficulty and mode
    const difficulty = getEnvValue('DIFFICULTY');
    const mode = getEnvValue('MODE');
    const memory = getEnvValue('MEMORY') || getEnvValue('MAX_MEMORY');
    const motd = getEnvValue('MOTD');
    
    if (difficulty) document.getElementById('setting-difficulty').value = difficulty;
    if (mode) document.getElementById('setting-mode').value = mode;
    if (memory) document.getElementById('setting-memory').value = memory;
    if (motd) document.getElementById('setting-motd').value = motd;
    
    // Load feature toggles
    const autopause = getEnvValue('ENABLE_AUTOPAUSE');
    const autostop = getEnvValue('ENABLE_AUTOSTOP');
    const whitelist = getEnvValue('ENABLE_WHITELIST');
    
    document.getElementById('setting-autopause').checked = autopause === 'true' || autopause === 'TRUE';
    document.getElementById('setting-autostop').checked = autostop === 'true' || autostop === 'TRUE';
    document.getElementById('setting-whitelist').checked = whitelist === 'true' || whitelist === 'TRUE';
    
  } catch(e) {
    console.log('Could not load settings:', e);
  }
}

async function applyServerSettings() {
  const difficulty = document.getElementById('setting-difficulty').value;
  const mode = document.getElementById('setting-mode').value;
  const memory = document.getElementById('setting-memory').value.trim();
  
  if (!memory) {
    toast('Please enter a memory value', 'err');
    return;
  }
  
  try {
    const configRes = await fetch(API + '/api/config');
    const config = await configRes.json();
    let content = config.content || '';
    
    // Update all three settings
    const settings = {
      DIFFICULTY: difficulty,
      MODE: mode,
      MAX_MEMORY: memory
    };
    
    Object.entries(settings).forEach(([key, value]) => {
      const pattern = new RegExp('^(\\s+' + key + ':\\s*)["\']?[^"\'\\n]+["\']?', 'im');
      
      if (pattern.test(content)) {
        content = content.replace(pattern, '      ' + key + ': "' + value + '"');
      } else {
        const envMatch = content.match(/environment:\n/);
        if (envMatch) {
          const envStart = content.indexOf(envMatch[0]) + envMatch[0].length;
          content = content.slice(0, envStart) + '      ' + key + ': "' + value + '"\n' + content.slice(envStart);
        }
      }
    });
    
    const res = await fetch(API + '/api/config', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content: content })
    });
    
    const data = await res.json();
    if (data.ok) {
      toast('Settings applied!', 'ok');
      loadSettings();
    } else {
      toast(data.message || 'Failed to apply settings', 'err');
    }
  } catch(e) {
    toast('Error applying settings', 'err');
  }
}

async function applyServerSettings() {
  const difficulty = document.getElementById('setting-difficulty').value;
  const mode = document.getElementById('setting-mode').value;
  const memory = document.getElementById('setting-memory').value.trim();
  
  if (!memory) {
    toast('Please enter a memory value', 'err');
    return;
  }
  
  try {
    const configRes = await fetch(API + '/api/config');
    const config = await configRes.json();
    let content = config.content || '';
    
    // Remove duplicate MAX_MEMORY if it exists
    content = content.replace(/MAX_MEMORY:.*?\n/g, '');
    content = content.replace(/MEMORY:.*?\n/g, '');
    
    // Update settings
    const settings = {
      DIFFICULTY: difficulty,
      MODE: mode,
      MAX_MEMORY: memory
    };
    
    Object.entries(settings).forEach(([key, value]) => {
      const pattern = new RegExp('^\\s+' + key + ':.*$', 'im');
      
      if (pattern.test(content)) {
        content = content.replace(pattern, '      ' + key + ': "' + value + '"');
      } else {
        // Find the last line of environment section and add before it
        const envMatch = content.match(/environment:\n([\s\S]*?)\n\s{0,4}[a-z]/);
        if (envMatch) {
          const lastEnvLine = envMatch[0].lastIndexOf('\n');
          const insertPoint = content.indexOf(envMatch[0]) + lastEnvLine;
          content = content.slice(0, insertPoint) + '\n      ' + key + ': "' + value + '"' + content.slice(insertPoint);
        }
      }
    });
    
    const res = await fetch(API + '/api/config', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content: content })
    });
    
    const data = await res.json();
    if (data.ok) {
      toast('Settings applied!', 'ok');
      loadSettings();
    } else {
      toast(data.message || 'Failed to apply settings', 'err');
    }
  } catch(e) {
    toast('Error applying settings', 'err');
  }
}

async function saveFeatures() {
  const features = {
    ENABLE_AUTOPAUSE: document.getElementById('setting-autopause').checked ? 'true' : 'false',
    ENABLE_AUTOSTOP: document.getElementById('setting-autostop').checked ? 'true' : 'false',
    ENABLE_WHITELIST: document.getElementById('setting-whitelist').checked ? 'true' : 'false'
  };
  
  try {
    const configRes = await fetch(API + '/api/config');
    const config = await configRes.json();
    let content = config.content || '';
    
    Object.entries(features).forEach(([key, value]) => {
      const pattern = new RegExp('(\\s+' + key + ':\\s*["\']?)([^"\'\\n]+)(["\']?)', 'i');
      
      if (pattern.test(content)) {
        // Replace existing variable
        content = content.replace(pattern, '$1' + value + '$3');
      } else {
        // Add new variable in the environment section
        const envMatch = content.match(/environment:\n/i);
        if (envMatch) {
          const insertPoint = content.indexOf(envMatch[0]) + envMatch[0].length;
          content = content.slice(0, insertPoint) + '      ' + key + ': ' + value + '\n' + content.slice(insertPoint);
        }
      }
    });
    
    const res = await fetch(API + '/api/config', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content: content })
    });
    
    const data = await res.json();
    if (data.ok) {
      toast('Features updated!', 'ok');
    } else {
      toast(data.message || 'Failed to update features', 'err');
    }
  } catch(e) {
    toast('Error updating features', 'err');
  }
}

// Settings tab switching
function goSettingsTab(tab) {
  document.querySelectorAll('.settings-tab-btn').forEach(function(b) { b.classList.remove('active'); });
  document.querySelectorAll('.settings-content').forEach(function(c) { c.style.display = 'none'; });
  
  var btn = document.querySelector('[onclick="goSettingsTab(\'' + tab + '\')"]');
  if (btn) btn.classList.add('active');
  
  var content = document.getElementById('settings-' + tab);
  if (content) content.style.display = 'block';
  
  // Force reflow to update colors
  void content.offsetHeight;
}

// Load server.properties file
function loadServerProperties() {
}

function applyResources() {
  var overrides = {};
  var hasAny = false;
  
  document.querySelectorAll('.res-checkbox').forEach(function(checkbox) {
    var key = checkbox.dataset.key;
    var input = document.querySelector('.prop-value[data-key="' + key + '"]');
    var value = input.value || input.placeholder;
    
    if (checkbox.checked) {
      hasAny = true;
      overrides[key] = value;
    }
  });
  
  if (!hasAny) {
    toast('Select at least one resource to override', 'err');
    return;
  }
  
  fetch(API + '/api/config')
    .then(function(r) { return r.json(); })
    .then(function(config) {
      var content = config.content || '';
      
      // Remove existing resources section
      content = content.replace(/\s+resources:[\s\S]*?(?=\n\s{0,4}[a-z]|\n\n|$)/g, '');
      
      if (hasAny) {
        var portsIndex = content.indexOf('ports:');
        if (portsIndex !== -1) {
          var resourcesSection = '\n    resources:\n';
          resourcesSection += '      limits:\n';
          if (overrides['cpu-limit']) resourcesSection += '        cpus: \'' + overrides['cpu-limit'] + '\'\n';
          if (overrides['memory-limit']) resourcesSection += '        memory: ' + overrides['memory-limit'] + '\n';
          resourcesSection += '      reservations:\n';
          if (overrides['cpu-reservation']) resourcesSection += '        cpus: \'' + overrides['cpu-reservation'] + '\'\n';
          if (overrides['memory-reservation']) resourcesSection += '        memory: ' + overrides['memory-reservation'] + '\n';
          
          content = content.slice(0, portsIndex) + resourcesSection + '    ' + content.slice(portsIndex);
        }
      }
      
      return fetch(API + '/api/config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content: content })
      });
    })
    .then(function(r) { return r.json(); })
    .then(function(data) {
      if (data.ok) {
        toast('Resources applied!', 'ok');
      } else {
        toast(data.message || 'Failed to apply resources', 'err');
      }
    })
    .catch(function(e) { toast('Error applying resources', 'err'); });
}

// Apply server properties overrides
function applyProperties() {
  var overrides = [];
  document.querySelectorAll('.prop-checkbox').forEach(function(checkbox) {
    if (checkbox.checked) {
      var key = checkbox.dataset.key;
      var input = document.querySelector('.prop-value[data-key="' + key + '"]');
      var value = input.value || input.placeholder;
      overrides.push(key + '=' + value);
    }
  });
  
  if (overrides.length === 0) {
    toast('Select at least one property to override', 'err');
    return;
  }
  
  fetch(API + '/api/config')
    .then(function(r) { return r.json(); })
    .then(function(config) {
      var content = config.content || '';
      var envVars = '';
      overrides.forEach(function(override) {
        var parts = override.split('=');
        var key = parts[0];
        var value = parts.slice(1).join('=');
        envVars += '      ' + key + ': "' + value + '"\n';
      });
      var envMatch = content.match(/environment:\n/);
      if (envMatch) {
        var envStart = content.indexOf(envMatch[0]) + envMatch[0].length;
        content = content.slice(0, envStart) + envVars + content.slice(envStart);
      }
      return fetch(API + '/api/config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content: content })
      });
    })
    .then(function(r) { return r.json(); })
    .then(function(data) {
      if (data.ok) {
        toast('Properties applied!', 'ok');
      } else {
        toast(data.message || 'Failed to apply properties', 'err');
      }
    })
    .catch(function(e) { toast('Error applying properties', 'err'); });
}

// Load settings when tab is opened
var settingsLoaded = false;
var originalGoTab = window.goTab;

window.goTab = function(name) {
  // Call original goTab
  document.querySelectorAll('.tab-panel').forEach(function(p) { p.classList.remove('active'); });
  document.querySelectorAll('.tab-btn').forEach(function(b) { b.classList.remove('active'); });
  var panel = document.getElementById('panel-' + name);
  var btn = document.getElementById('tab-' + name);
  if (panel) panel.classList.add('active');
  if (btn) btn.classList.add('active');
  
  // Lazy-load tab data
  if (!tabLoaded[name]) {
    tabLoaded[name] = true;
    if (name === 'files') loadFiles('/data');
    if (name === 'mods') fetchMods();
    if (name === 'config') fetchConfig();
  }
  
  // Settings-specific logic
  if (name === 'settings' && !settingsLoaded) {
    settingsLoaded = true;
    loadSettings();
  }
  if (name === 'config') {
    fetchConfig();
  }
};

</script>
</body>
</html>`
