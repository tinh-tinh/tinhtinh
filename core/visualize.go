package core

import (
	"encoding/json"
	"html/template"
	"os"
)

// ---------------------------------------------------------------------------
// Tree node types
// ---------------------------------------------------------------------------

// RouteNode describes a single HTTP route registered by a controller.
type RouteNode struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// ControllerNode represents a controller and its registered routes.
type ControllerNode struct {
	Name   string      `json:"name"`
	Routes []RouteNode `json:"routes"`
}

// ProviderNode represents a provider registered in a module.
type ProviderNode struct {
	Name   string        `json:"name"`
	Scope  Scope         `json:"scope"`
	Status ProvideStatus `json:"status"`
}

// ModuleNode is a recursive description of the module hierarchy.
type ModuleNode struct {
	Name        string           `json:"name"`
	Scope       Scope            `json:"scope"`
	Controllers []ControllerNode `json:"controllers,omitempty"`
	Providers   []ProviderNode   `json:"providers,omitempty"`
	Imports     []*ModuleNode    `json:"imports,omitempty"`
}

// ---------------------------------------------------------------------------
// Tree builder
// ---------------------------------------------------------------------------

// buildModuleNode converts a live DynamicModule into a ModuleNode tree.
// It groups Router entries by controller name to reconstruct controllers,
// and reads DataProviders that are PRIVATE (own providers, not re-exported
// imports) directly.
func buildModuleNode(module *DynamicModule) *ModuleNode {
	node := &ModuleNode{
		Name:  module.Name,
		Scope: module.Scope,
	}

	// Use SnapshotRouters when available (set before free() clears Routers),
	// otherwise fall back to Routers (root module never calls free).
	routers := module.SnapshotRouters
	if len(routers) == 0 {
		routers = module.Routers
	}

	// --- Collect own controllers from Routers ---
	ctrlMap := make(map[string]*ControllerNode)
	ctrlOrder := []string{}
	for _, r := range routers {
		if _, exists := ctrlMap[r.Name]; !exists {
			ctrlMap[r.Name] = &ControllerNode{Name: r.Name}
			ctrlOrder = append(ctrlOrder, r.Name)
		}
		rn := RouteNode{Method: r.Method, Path: r.Path}
		ctrlMap[r.Name].Routes = append(ctrlMap[r.Name].Routes, rn)
	}
	for _, name := range ctrlOrder {
		node.Controllers = append(node.Controllers, *ctrlMap[name])
	}

	// --- Collect own providers ---
	for _, p := range module.DataProviders {
		node.Providers = append(node.Providers, ProviderNode{
			Name:   string(p.GetName()),
			Scope:  p.GetScope(),
			Status: p.GetStatus(),
		})
	}

	// --- Recurse into sub-modules ---
	for _, sub := range module.SubModules {
		node.Imports = append(node.Imports, buildModuleNode(sub))
	}

	return node
}

// ---------------------------------------------------------------------------
// App methods
// ---------------------------------------------------------------------------

// GetTree builds and returns the full module/controller/provider tree rooted
// at the App's module. Call this after CreateFactory.
func (app *App) GetTree() *ModuleNode {
	dynMod, ok := app.Module.(*DynamicModule)
	if !ok {
		return nil
	}
	return buildModuleNode(dynMod)
}

// ---------------------------------------------------------------------------
// HTML generation
// ---------------------------------------------------------------------------

// Visualize generates a self-contained, interactive HTML file that renders
// the module/controller/provider tree using D3.js.
// outputPath is the file path where the HTML should be written.
func (app *App) Visualize(outputPath string) error {
	tree := app.GetTree()
	if tree == nil {
		return nil
	}

	treeJSON, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		return err
	}

	tmpl, err := template.New("visualize").Parse(visualizeTemplate)
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	data := struct {
		TreeJSON template.JS
	}{
		TreeJSON: template.JS(treeJSON),
	}

	return tmpl.Execute(f, data)
}

// ---------------------------------------------------------------------------
// HTML template (self-contained, uses D3 v7 from CDN)
// ---------------------------------------------------------------------------

const visualizeTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
<title>TinhTinh Module Tree</title>
<style>
  :root {
    --bg: #0f1117;
    --surface: #1a1d27;
    --surface2: #22263a;
    --border: #2e3352;
    --text: #e2e8f0;
    --text-muted: #8892a4;
    --module-color: #6366f1;
    --controller-color: #22d3ee;
    --provider-color: #34d399;
    --route-color: #f59e0b;
    --export-color: #a78bfa;
    --link-color: #3b4299;
    --font: 'Inter', 'Segoe UI', system-ui, sans-serif;
  }

  * { box-sizing: border-box; margin: 0; padding: 0; }

  body {
    font-family: var(--font);
    background: var(--bg);
    color: var(--text);
    min-height: 100vh;
    overflow: hidden;
  }

  header {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px 24px;
    background: var(--surface);
    border-bottom: 1px solid var(--border);
    position: fixed;
    top: 0; left: 0; right: 0;
    z-index: 100;
  }

  header .logo {
    font-size: 20px;
    font-weight: 800;
    background: linear-gradient(135deg, #6366f1, #22d3ee);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    letter-spacing: -0.5px;
  }

  header .subtitle {
    font-size: 13px;
    color: var(--text-muted);
  }

  .legend {
    display: flex;
    gap: 16px;
    margin-left: auto;
    align-items: center;
    flex-wrap: wrap;
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    color: var(--text-muted);
  }

  .legend-dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    flex-shrink: 0;
  }

  #canvas {
    position: fixed;
    top: 57px; left: 0; right: 0; bottom: 0;
    overflow: hidden;
  }

  svg { width: 100%; height: 100%; }

  /* Tree links */
  .link {
    fill: none;
    stroke: var(--link-color);
    stroke-width: 1.5px;
    stroke-opacity: 0.7;
  }

  /* Nodes */
  .node circle {
    stroke-width: 2.5px;
    cursor: pointer;
    transition: filter 0.2s;
  }

  .node circle:hover { filter: brightness(1.3); }

  .node text {
    font-family: var(--font);
    font-size: 13px;
    fill: var(--text);
    pointer-events: none;
  }

  .node .label-bg {
    fill: var(--surface);
    rx: 4px;
    opacity: 0.85;
  }

  /* Tooltip */
  #tooltip {
    position: fixed;
    background: var(--surface2);
    border: 1px solid var(--border);
    border-radius: 10px;
    padding: 14px 18px;
    font-size: 13px;
    pointer-events: none;
    opacity: 0;
    transition: opacity 0.2s;
    max-width: 320px;
    z-index: 200;
    box-shadow: 0 8px 32px rgba(0,0,0,0.5);
  }

  #tooltip.visible { opacity: 1; }

  #tooltip h3 {
    font-size: 14px;
    margin-bottom: 10px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  #tooltip .badge {
    font-size: 10px;
    font-weight: 700;
    padding: 2px 7px;
    border-radius: 99px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .badge-module   { background: rgba(99,102,241,0.25); color: #818cf8; }
  .badge-controller { background: rgba(34,211,238,0.2); color: #22d3ee; }
  .badge-provider { background: rgba(52,211,153,0.2); color: #34d399; }

  #tooltip .section-title {
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--text-muted);
    margin: 8px 0 4px;
  }

  #tooltip ul {
    list-style: none;
    display: flex;
    flex-direction: column;
    gap: 3px;
  }

  #tooltip li {
    font-size: 12px;
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .method-badge {
    font-size: 10px;
    font-weight: 700;
    padding: 1px 5px;
    border-radius: 4px;
    text-transform: uppercase;
    min-width: 40px;
    text-align: center;
  }

  .GET    { background: rgba(52,211,153,0.2); color: #34d399; }
  .POST   { background: rgba(251,191,36,0.2); color: #fbbf24; }
  .PUT    { background: rgba(96,165,250,0.2); color: #60a5fa; }
  .PATCH  { background: rgba(167,139,250,0.2); color: #a78bfa; }
  .DELETE { background: rgba(239,68,68,0.2); color: #f87171; }

  .controls {
    position: fixed;
    bottom: 24px;
    right: 24px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    z-index: 100;
  }

  .ctrl-btn {
    width: 38px;
    height: 38px;
    background: var(--surface2);
    border: 1px solid var(--border);
    border-radius: 8px;
    color: var(--text);
    font-size: 18px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.15s, transform 0.1s;
  }

  .ctrl-btn:hover { background: var(--border); transform: scale(1.05); }
  .ctrl-btn:active { transform: scale(0.95); }

  /* Expand all / Collapse all */
  .toggle-bar {
    position: fixed;
    bottom: 24px;
    left: 24px;
    display: flex;
    gap: 8px;
    z-index: 100;
  }

  .toggle-btn {
    background: var(--surface2);
    border: 1px solid var(--border);
    border-radius: 8px;
    color: var(--text);
    font-size: 12px;
    padding: 8px 14px;
    cursor: pointer;
    transition: background 0.15s;
  }

  .toggle-btn:hover { background: var(--border); }
</style>
</head>
<body>

<header>
  <span class="logo">TinhTinh</span>
  <span class="subtitle">Module Tree Visualizer</span>
  <div class="legend">
    <div class="legend-item"><div class="legend-dot" style="background:#6366f1"></div> Module</div>
    <div class="legend-item"><div class="legend-dot" style="background:#22d3ee"></div> Controller</div>
    <div class="legend-item"><div class="legend-dot" style="background:#34d399"></div> Provider (private)</div>
    <div class="legend-item"><div class="legend-dot" style="background:#a78bfa"></div> Provider (public)</div>
  </div>
</header>

<div id="canvas"><svg id="tree-svg"><g id="root-g"></g></svg></div>

<div id="tooltip"></div>

<div class="controls">
  <button class="ctrl-btn" id="zoom-in"  title="Zoom in">+</button>
  <button class="ctrl-btn" id="zoom-out" title="Zoom out">−</button>
  <button class="ctrl-btn" id="zoom-fit" title="Fit" style="font-size:14px">⊡</button>
</div>

<div class="toggle-bar">
  <button class="toggle-btn" id="expand-all">Expand All</button>
  <button class="toggle-btn" id="collapse-all">Collapse All</button>
</div>

<script src="https://cdn.jsdelivr.net/npm/d3@7/dist/d3.min.js"></script>
<script>
// ── Data ──────────────────────────────────────────────────────────────────
const RAW = {{.TreeJSON}};

// ── Flatten the module tree into a D3-friendly hierarchy ──────────────────
//
// Each module node may have:
//   - children that are sub-modules (imports)
//   - synthetic leaf children for controllers (and their routes)
//   - synthetic leaf children for providers

function buildHierarchy(mod, depth) {
  const node = {
    id: Math.random().toString(36).slice(2),
    name: mod.name || 'AppModule',
    kind: 'module',
    scope: mod.scope,
    _controllers: mod.controllers || [],
    _providers: mod.providers || [],
    children: [],
    _collapsed: depth > 0,
  };

  // Controller children
  (mod.controllers || []).forEach(ctrl => {
    node.children.push({
      id: Math.random().toString(36).slice(2),
      name: ctrl.name,
      kind: 'controller',
      _routes: ctrl.routes || [],
      children: [],
      _collapsed: false,
    });
  });

  // Provider children
  (mod.providers || []).forEach(prov => {
    node.children.push({
      id: Math.random().toString(36).slice(2),
      name: prov.name,
      kind: 'provider',
      scope: prov.scope,
      status: prov.status,
      children: [],
      _collapsed: false,
    });
  });

  // Sub-module children
  (mod.imports || []).forEach(sub => {
    node.children.push(buildHierarchy(sub, depth + 1));
  });

  return node;
}

const treeData = buildHierarchy(RAW, 0);

// ── Layout constants ───────────────────────────────────────────────────────
const NODE_SEP_X = 260;
const NODE_SEP_Y = 52;

// ── Colors ────────────────────────────────────────────────────────────────
const COLOR = {
  module:     '#6366f1',
  controller: '#22d3ee',
  provider_private: '#34d399',
  provider_public:  '#a78bfa',
};

function nodeColor(d) {
  if (d.data.kind === 'module')     return COLOR.module;
  if (d.data.kind === 'controller') return COLOR.controller;
  if (d.data.kind === 'provider')
    return d.data.status === 'public' ? COLOR.provider_public : COLOR.provider_private;
  return '#888';
}

// ── SVG / zoom setup ──────────────────────────────────────────────────────
const svg  = d3.select('#tree-svg');
const rootG = d3.select('#root-g');

const zoom = d3.zoom()
  .scaleExtent([0.1, 3])
  .on('zoom', e => rootG.attr('transform', e.transform));

svg.call(zoom).on('dblclick.zoom', null);

// ── Tooltip ───────────────────────────────────────────────────────────────
const tooltip = document.getElementById('tooltip');

function showTooltip(event, d) {
  const data = d.data;
  let html = '';

  if (data.kind === 'module') {
    html = '<h3>' + escHtml(data.name) + ' <span class="badge badge-module">Module</span></h3>';
    html += '<div class="section-title">Scope</div>';
    html += '<ul><li>' + escHtml(data.scope || 'global') + '</li></ul>';
    if (data._controllers && data._controllers.length) {
      html += '<div class="section-title">Controllers (' + data._controllers.length + ')</div><ul>';
      data._controllers.forEach(c => {
        html += '<li>🎮 ' + escHtml(c.name) + '</li>';
        (c.routes || []).forEach(r => {
          html += '<li style="padding-left:12px"><span class="method-badge ' + r.method + '">' + r.method + '</span> ' + escHtml(r.path || '/') + '</li>';
        });
      });
      html += '</ul>';
    }
    if (data._providers && data._providers.length) {
      html += '<div class="section-title">Providers (' + data._providers.length + ')</div><ul>';
      data._providers.forEach(p => {
        const badge = p.status === 'public' ? '🟣' : '🟢';
        html += '<li>' + badge + ' ' + escHtml(p.name) + ' <span style="color:var(--text-muted);font-size:11px">' + (p.scope || '') + '</span></li>';
      });
      html += '</ul>';
    }
  } else if (data.kind === 'controller') {
    html = '<h3>🎮 ' + escHtml(data.name) + ' <span class="badge badge-controller">Controller</span></h3>';
    if (data._routes && data._routes.length) {
      html += '<div class="section-title">Routes</div><ul>';
      data._routes.forEach(r => {
        html += '<li><span class="method-badge ' + r.method + '">' + r.method + '</span> <span style="font-family:monospace">' + escHtml(r.path || '/') + '</span></li>';
      });
      html += '</ul>';
    } else {
      html += '<div style="color:var(--text-muted);font-size:12px;margin-top:6px">No routes registered</div>';
    }
  } else if (data.kind === 'provider') {
    const icon = data.status === 'public' ? '🟣' : '🟢';
    html = '<h3>' + icon + ' ' + escHtml(data.name) + ' <span class="badge badge-provider">Provider</span></h3>';
    html += '<div class="section-title">Details</div><ul>';
    html += '<li>Scope: ' + escHtml(data.scope || 'global') + '</li>';
    html += '<li>Status: ' + escHtml(data.status || 'private') + '</li>';
    html += '</ul>';
  }

  tooltip.innerHTML = html;
  tooltip.classList.add('visible');
  moveTooltip(event);
}

function moveTooltip(event) {
  const x = event.clientX + 16;
  const y = event.clientY + 16;
  const w = tooltip.offsetWidth;
  const h = tooltip.offsetHeight;
  const winW = window.innerWidth;
  const winH = window.innerHeight;
  tooltip.style.left = Math.min(x, winW - w - 16) + 'px';
  tooltip.style.top  = Math.min(y, winH - h - 16) + 'px';
}

function hideTooltip() {
  tooltip.classList.remove('visible');
}

function escHtml(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');
}

// ── D3 tree rendering ─────────────────────────────────────────────────────
let root;

function update(source) {
  const treeLayout = d3.tree()
    .nodeSize([NODE_SEP_Y, NODE_SEP_X])
    .separation((a, b) => a.parent === b.parent ? 1.2 : 1.6);

  // Build fresh hierarchy each render
  root = d3.hierarchy(treeData, d => d._collapsed ? null : (d.children.length ? d.children : null));
  treeLayout(root);

  const duration = 300;

  // --- Links ---
  const linkGen = d3.linkHorizontal()
    .x(d => d.y)
    .y(d => d.x);

  const links = rootG.selectAll('.link')
    .data(root.links(), d => d.target.data.id);

  links.enter()
    .append('path')
    .attr('class', 'link')
    .attr('d', d => {
      const o = { x: source ? source.x0 : 0, y: source ? source.y0 : 0 };
      return linkGen({ source: o, target: o });
    })
    .merge(links)
    .transition().duration(duration)
    .attr('d', linkGen);

  links.exit()
    .transition().duration(duration)
    .attr('d', d => {
      const o = { x: source ? source.x : 0, y: source ? source.y : 0 };
      return linkGen({ source: o, target: o });
    })
    .remove();

  // --- Nodes ---
  const nodes = rootG.selectAll('.node')
    .data(root.descendants(), d => d.data.id);

  const nodeEnter = nodes.enter()
    .append('g')
    .attr('class', 'node')
    .attr('transform', d => 'translate(' + (source ? source.y0 : 0) + ',' + (source ? source.x0 : 0) + ')')
    .on('click', (event, d) => {
      event.stopPropagation();
      d.data._collapsed = !d.data._collapsed;
      update(d);
    })
    .on('mouseenter', showTooltip)
    .on('mousemove', moveTooltip)
    .on('mouseleave', hideTooltip);

  // Background rect for label
  nodeEnter.append('rect')
    .attr('class', 'label-bg')
    .attr('rx', 4)
    .attr('ry', 4)
    .attr('fill', '#1a1d27')
    .attr('opacity', 0.85);

  // Circle
  nodeEnter.append('circle')
    .attr('r', 0)
    .style('fill', nodeColor)
    .style('stroke', nodeColor);

  // Label text
  nodeEnter.append('text')
    .attr('dy', '.35em')
    .attr('text-anchor', 'start')
    .text(d => {
      const icons = { module: '▣', controller: '⬡', provider: '◆' };
      return (icons[d.data.kind] || '') + ' ' + d.data.name;
    })
    .style('fill', nodeColor);

  const nodeUpdate = nodeEnter.merge(nodes)
    .transition().duration(duration)
    .attr('transform', d => 'translate(' + d.y + ',' + d.x + ')');

  nodeUpdate.select('circle')
    .attr('r', d => d.data.kind === 'module' ? 10 : 7)
    .style('fill', d => d.data._collapsed ? nodeColor(d) : 'rgba(0,0,0,0)');

  // Position and size the label background
  nodeEnter.each(function(d) {
    const g = d3.select(this);
    const text = g.select('text');
    setTimeout(() => {
      try {
        const bb = text.node().getBBox();
        g.select('.label-bg')
          .attr('x', bb.x - 4)
          .attr('y', bb.y - 2)
          .attr('width', bb.width + 8)
          .attr('height', bb.height + 4);
      } catch (e) {}
    }, duration + 50);
  });

  nodeUpdate.select('text')
    .attr('x', d => (d.data.kind === 'module' ? 16 : 12));

  nodes.exit()
    .transition().duration(duration)
    .attr('transform', d => 'translate(' + (source ? source.y : 0) + ',' + (source ? source.x : 0) + ')')
    .remove();

  // Save positions for transition origins
  root.descendants().forEach(d => {
    d.x0 = d.x;
    d.y0 = d.y;
  });
}

// ── Initial render + fit ───────────────────────────────────────────────────
update(null);

function fitView() {
  const bounds = rootG.node().getBBox();
  const canvasEl = document.getElementById('canvas');
  const w = canvasEl.clientWidth;
  const h = canvasEl.clientHeight;
  if (bounds.width === 0 || bounds.height === 0) return;
  const scale = Math.min(w / (bounds.width + 80), h / (bounds.height + 80), 1.5) * 0.9;
  const tx = w / 2 - scale * (bounds.x + bounds.width / 2);
  const ty = h / 2 - scale * (bounds.y + bounds.height / 2);
  svg.transition().duration(600)
    .call(zoom.transform, d3.zoomIdentity.translate(tx, ty).scale(scale));
}

setTimeout(fitView, 400);

// ── Controls ──────────────────────────────────────────────────────────────
document.getElementById('zoom-in').addEventListener('click', () => svg.transition().call(zoom.scaleBy, 1.3));
document.getElementById('zoom-out').addEventListener('click', () => svg.transition().call(zoom.scaleBy, 0.77));
document.getElementById('zoom-fit').addEventListener('click', fitView);

function setAllCollapsed(val, node) {
  node._collapsed = val;
  (node.children || []).forEach(c => setAllCollapsed(val, c));
}

document.getElementById('expand-all').addEventListener('click', () => {
  setAllCollapsed(false, treeData);
  update(null);
  setTimeout(fitView, 400);
});

document.getElementById('collapse-all').addEventListener('click', () => {
  setAllCollapsed(true, treeData);
  treeData._collapsed = false; // keep root visible
  update(null);
  setTimeout(fitView, 400);
});
</script>
</body>
</html>
`
