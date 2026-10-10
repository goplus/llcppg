// Outline sidebar behaviour for `llcppg -doc`.
//
// internal/godoc renders the documentation outline as a nested <ul> tree using
// pkg.go.dev's class names and #anchors. This script enhances that static tree
// with the behaviour a reader expects from pkg.go.dev's sidebar, without any
// build step or network access.
//
// It is a faithful vanilla-JS port of pkgsite's
// static/shared/outline/tree.ts (BSD-licensed, Copyright The Go Authors), which
// implements the WAI-ARIA Treeview pattern:
//
//   - Collapse/expand: every item that has a child <ul> is collapsible and
//     starts collapsed (aria-expanded="false"); clicking it (or Enter/Space,
//     ArrowRight/ArrowLeft) toggles its group. The CSS in page.css (also copied
//     from pkgsite's tree.css) draws the triangle toggles and hides collapsed
//     groups via `a[aria-expanded='true'] + ul[role='group']`.
//   - Scroll-spy: an IntersectionObserver watches the body section each item
//     targets; the item for the section in view is selected (aria-selected) and
//     its ancestors are expanded so it is revealed.
//   - Keyboard navigation: arrows/Home/End/typeahead move focus between visible
//     items with a roving tabindex, matching pkgsite.
//
// It also fills the mobile <select> fallback from the same tree so narrow
// screens get a jump menu instead of the full sidebar.
//
// All anchors come from the server-rendered tree, so the markup and the body
// agree by construction; this script only reads ids, it never invents them.
(function () {
  "use strict";

  var tree = document.querySelector(".js-tree");
  if (!tree) {
    return;
  }

  // debounce delays fn until wait ms have passed without another call, matching
  // pkgsite's scroll-spy debouncing so rapid intersection events settle once.
  function debounce(fn, wait) {
    var timeout = null;
    return function () {
      var args = arguments;
      var later = function () {
        timeout = null;
        fn.apply(null, args);
      };
      if (timeout) {
        clearTimeout(timeout);
      }
      timeout = setTimeout(later, wait);
    };
  }

  // slug mirrors pkgsite's group-id generation: non-word runs collapse to "_".
  function slug(s) {
    return s.replace(/[\W_]+/g, "_");
  }

  // -- TreeItem -----------------------------------------------------------

  // TreeItem wraps one <a>/<span> in the tree, wiring the ARIA attributes and
  // its own event handlers. An item is expandable when it has a sibling <ul>;
  // that group is marked role="group" and linked via aria-owns, and the item
  // starts collapsed.
  function TreeItem(el, controller, group) {
    this.el = el;
    this.controller = controller;
    this.groupTreeitem = group;
    this.label = (el.textContent || "").trim();
    this.depth = (group ? group.depth : 0) + 1;
    this.index = 0;
    this.isExpandable = false;
    this.isVisible = false;
    this.isInGroup = !!group;

    el.tabIndex = -1;

    var parent = el.parentElement;
    if (parent && parent.tagName.toLowerCase() === "li") {
      parent.setAttribute("role", "none");
    }
    el.setAttribute("aria-level", String(this.depth));
    var ariaLabel = el.getAttribute("aria-label");
    if (ariaLabel) {
      this.label = ariaLabel.trim();
    }

    // An item becomes expandable when a sibling <ul> follows it. Collapsed by
    // default (aria-expanded="false"), so the sidebar opens compact.
    var curr = el.nextElementSibling;
    while (curr) {
      if (curr.tagName.toLowerCase() === "ul") {
        var groupId = slug((group ? group.label : "") + " nav group " + this.label);
        el.setAttribute("aria-owns", groupId);
        el.setAttribute("aria-expanded", "false");
        curr.setAttribute("role", "group");
        curr.setAttribute("id", groupId);
        this.isExpandable = true;
        break;
      }
      curr = curr.nextElementSibling;
    }

    if (!el.getAttribute("role")) {
      el.setAttribute("role", "treeitem");
    }

    var self = this;
    el.addEventListener("keydown", function (e) {
      self.handleKeydown(e);
    });
    el.addEventListener("click", function (e) {
      self.handleClick(e);
    });
    el.addEventListener("focus", function () {
      self.handleFocus();
    });
    el.addEventListener("blur", function () {
      self.handleBlur();
    });
  }

  TreeItem.prototype.isExpanded = function () {
    return this.isExpandable && this.el.getAttribute("aria-expanded") === "true";
  };

  TreeItem.prototype.isSelected = function () {
    return this.el.getAttribute("aria-selected") === "true";
  };

  TreeItem.prototype.handleClick = function (event) {
    // Only act on clicks on this item itself (or its first child), so a click
    // on a nested item does not toggle its ancestors.
    if (event.target !== this.el && event.target !== this.el.firstElementChild) {
      return;
    }
    if (this.isExpandable) {
      // Toggle purely on the expanded state. pkgsite also requires the item to
      // be selected before collapsing, but llcppg's scroll-spy re-selects the
      // in-view section continuously, which would deselect this item between
      // clicks and make a second click re-expand instead of collapse. Keying the
      // toggle on isExpanded() alone makes click reliably flip the group.
      if (this.isExpanded()) {
        this.controller.collapseTreeitem(this);
      } else {
        this.controller.expandTreeitem(this);
      }
      event.stopPropagation();
    }
    this.controller.setSelected(this);
  };

  TreeItem.prototype.handleFocus = function () {
    var el = this.el;
    if (this.isExpandable && el.firstElementChild) {
      el = el.firstElementChild;
    }
    el.classList.add("focus");
  };

  TreeItem.prototype.handleBlur = function () {
    var el = this.el;
    if (this.isExpandable && el.firstElementChild) {
      el = el.firstElementChild;
    }
    el.classList.remove("focus");
  };

  TreeItem.prototype.handleKeydown = function (event) {
    if (event.altKey || event.ctrlKey || event.metaKey) {
      return;
    }
    var captured = false;
    switch (event.key) {
      case " ":
      case "Enter":
        if (this.isExpandable) {
          // Same reliable-toggle rule as handleClick: collapse whenever the
          // group is open, independent of the scroll-spy's selection.
          if (this.isExpanded()) {
            this.controller.collapseTreeitem(this);
          } else {
            this.controller.expandTreeitem(this);
          }
          captured = true;
        } else {
          event.stopPropagation();
        }
        this.controller.setSelected(this);
        break;
      case "ArrowUp":
        this.controller.setFocusToPreviousItem(this);
        captured = true;
        break;
      case "ArrowDown":
        this.controller.setFocusToNextItem(this);
        captured = true;
        break;
      case "ArrowRight":
        if (this.isExpandable) {
          if (this.isExpanded()) {
            this.controller.setFocusToNextItem(this);
          } else {
            this.controller.expandTreeitem(this);
          }
        }
        captured = true;
        break;
      case "ArrowLeft":
        if (this.isExpandable && this.isExpanded()) {
          this.controller.collapseTreeitem(this);
          captured = true;
        } else if (this.isInGroup) {
          this.controller.setFocusToParentItem(this);
          captured = true;
        }
        break;
      case "Home":
        this.controller.setFocusToFirstItem();
        captured = true;
        break;
      case "End":
        this.controller.setFocusToLastItem();
        captured = true;
        break;
      default:
        if (event.key.length === 1 && event.key.match(/\S/)) {
          if (event.key === "*") {
            this.controller.expandAllSiblingItems(this);
          } else {
            this.controller.setFocusByFirstCharacter(this, event.key);
          }
          captured = true;
        }
        break;
    }
    if (captured) {
      event.stopPropagation();
      event.preventDefault();
    }
  };

  // -- TreeNavController ---------------------------------------------------

  function TreeNavController(el) {
    this.el = el;
    this.treeitems = [];
    this.firstChars = [];
    this.firstTreeitem = null;
    this.lastTreeitem = null;
    this.observerCallbacks = [];
    this.init();
  }

  TreeNavController.prototype.init = function () {
    this.handleResize();
    var self = this;
    window.addEventListener("resize", function () {
      self.handleResize();
    });
    this.findTreeItems();
    // Open the top-level containers ("Documentation") once on load so their
    // section links (Overview, Index, Functions, Types, …) are visible, but
    // leave every deeper group collapsed. This is the "collapsed by default"
    // the sidebar wants: compact, with only the section headers showing until
    // the reader expands a group.
    this.expandTopLevelItems();
    this.updateVisibleTreeitems();
    this.observeTargets();
    if (this.firstTreeitem) {
      this.firstTreeitem.el.tabIndex = 0;
    }
  };

  // expandTopLevelItems expands only the depth-1 containers, revealing the
  // outline's section links while keeping the nested func/type groups collapsed.
  TreeNavController.prototype.expandTopLevelItems = function () {
    for (var i = 0; i < this.treeitems.length; i++) {
      var ti = this.treeitems[i];
      if (ti.isExpandable && ti.depth === 1) {
        ti.el.setAttribute("aria-expanded", "true");
      }
    }
  };

  TreeNavController.prototype.handleResize = function () {
    // Expose the tree's pixel height so the CSS can cap the scrollable group,
    // matching pkgsite's --js-tree-height.
    this.el.style.setProperty("--js-tree-height", "100vh");
    this.el.style.setProperty("--js-tree-height", this.el.clientHeight + "px");
  };

  TreeNavController.prototype.observeTargets = function () {
    var self = this;
    this.addObserver(function (treeitem) {
      // Scroll-spy reveals and selects the item for the section in view, but it
      // must not force the top-level (level 1) "Documentation" node open: doing
      // so on load expands the whole outline, defeating "collapsed by default".
      // So expand only the mid-level ancestors (type groups, etc.) and leave the
      // level-1 container in whatever state the reader last set.
      self.expandAncestorsBelowTop(treeitem);
      self.setSelected(treeitem);
    });

    if (!("IntersectionObserver" in window)) {
      return;
    }
    var targets = {};
    var order = [];
    var observer = new IntersectionObserver(
      function (entries) {
        for (var i = 0; i < entries.length; i++) {
          var entry = entries[i];
          var hit = entry.isIntersecting || entry.intersectionRatio === 1;
          if (!(entry.target.id in targets)) {
            order.push(entry.target.id);
          }
          targets[entry.target.id] = hit;
        }
        for (var j = 0; j < order.length; j++) {
          var id = order[j];
          if (targets[id]) {
            var active = self.findByHrefId(id);
            if (active) {
              for (var k = 0; k < self.observerCallbacks.length; k++) {
                self.observerCallbacks[k](active);
              }
            }
            break;
          }
        }
      },
      { threshold: 1.0, rootMargin: "-60px 0px 0px 0px" }
    );
    for (var t = 0; t < this.treeitems.length; t++) {
      var href = this.treeitems[t].el.getAttribute("href");
      if (!href) {
        continue;
      }
      var id = decodeURIComponent(
        href.replace(window.location.origin, "").replace("/", "").replace("#", "")
      );
      var target = document.getElementById(id);
      if (target && id.indexOf("example-") !== 0) {
        observer.observe(target);
      }
    }
  };

  // findByHrefId returns the treeitem whose href is exactly "#<id>", where id is
  // the already-decoded target id the observer resolved. Both sides are compared
  // on the decoded href so a percent-encoded anchor still matches, and the
  // leading "#" keeps a short id (e.g. "Close") from matching a longer href
  // (e.g. "#DB.Close").
  TreeNavController.prototype.findByHrefId = function (id) {
    var needle = "#" + id;
    for (var i = 0; i < this.treeitems.length; i++) {
      var href = this.treeitems[i].el.getAttribute("href");
      if (!href) {
        continue;
      }
      var dec = decodeURIComponent(href);
      if (dec.length >= needle.length && dec.indexOf(needle) === dec.length - needle.length) {
        return this.treeitems[i];
      }
    }
    return null;
  };

  TreeNavController.prototype.addObserver = function (fn, delay) {
    this.observerCallbacks.push(debounce(fn, delay == null ? 200 : delay));
  };

  TreeNavController.prototype.setFocusToNextItem = function (currentItem) {
    var nextItem = null;
    for (var i = currentItem.index + 1; i < this.treeitems.length; i++) {
      if (this.treeitems[i].isVisible) {
        nextItem = this.treeitems[i];
        break;
      }
    }
    if (nextItem) {
      this.setFocusToItem(nextItem);
    }
  };

  TreeNavController.prototype.setFocusToPreviousItem = function (currentItem) {
    var prevItem = null;
    for (var i = currentItem.index - 1; i > -1; i--) {
      if (this.treeitems[i].isVisible) {
        prevItem = this.treeitems[i];
        break;
      }
    }
    if (prevItem) {
      this.setFocusToItem(prevItem);
    }
  };

  TreeNavController.prototype.setFocusToParentItem = function (currentItem) {
    if (currentItem.groupTreeitem) {
      this.setFocusToItem(currentItem.groupTreeitem);
    }
  };

  TreeNavController.prototype.setFocusToFirstItem = function () {
    if (this.firstTreeitem) {
      this.setFocusToItem(this.firstTreeitem);
    }
  };

  TreeNavController.prototype.setFocusToLastItem = function () {
    if (this.lastTreeitem) {
      this.setFocusToItem(this.lastTreeitem);
    }
  };

  TreeNavController.prototype.setSelected = function (currentItem) {
    var expanded = this.el.querySelectorAll('[aria-expanded="true"]');
    for (var i = 0; i < expanded.length; i++) {
      var l1 = expanded[i];
      if (l1 === currentItem.el) {
        continue;
      }
      var sib = l1.nextElementSibling;
      if (!(sib && sib.contains(currentItem.el))) {
        l1.setAttribute("aria-expanded", "false");
      }
    }
    var selected = this.el.querySelectorAll("[aria-selected]");
    for (var j = 0; j < selected.length; j++) {
      if (selected[j] !== currentItem.el) {
        selected[j].setAttribute("aria-selected", "false");
      }
    }
    currentItem.el.setAttribute("aria-selected", "true");
    this.updateVisibleTreeitems();
    this.setFocusToItem(currentItem, false);
  };

  TreeNavController.prototype.expandTreeitem = function (treeitem) {
    var currentItem = treeitem;
    while (currentItem) {
      if (currentItem.isExpandable) {
        currentItem.el.setAttribute("aria-expanded", "true");
      }
      currentItem = currentItem.groupTreeitem;
    }
    this.updateVisibleTreeitems();
  };

  // expandAncestorsBelowTop expands the item and its ancestors like
  // expandTreeitem, but skips any ancestor at tree depth 1 (the top-level
  // "Documentation" container). The scroll-spy uses it so that having a section
  // in view does not auto-open the whole outline on load; the sidebar stays
  // compact until the reader expands a group themselves.
  TreeNavController.prototype.expandAncestorsBelowTop = function (treeitem) {
    var currentItem = treeitem;
    while (currentItem) {
      if (currentItem.isExpandable && currentItem.depth > 1) {
        currentItem.el.setAttribute("aria-expanded", "true");
      }
      currentItem = currentItem.groupTreeitem;
    }
    this.updateVisibleTreeitems();
  };

  TreeNavController.prototype.expandAllSiblingItems = function (currentItem) {
    for (var i = 0; i < this.treeitems.length; i++) {
      var ti = this.treeitems[i];
      if (ti.groupTreeitem === currentItem.groupTreeitem && ti.isExpandable) {
        this.expandTreeitem(ti);
      }
    }
  };

  TreeNavController.prototype.collapseTreeitem = function (currentItem) {
    var groupTreeitem = currentItem.isExpanded() ? currentItem : currentItem.groupTreeitem;
    if (groupTreeitem) {
      groupTreeitem.el.setAttribute("aria-expanded", "false");
      this.updateVisibleTreeitems();
      this.setFocusToItem(groupTreeitem);
    }
  };

  TreeNavController.prototype.setFocusByFirstCharacter = function (currentItem, char) {
    char = char.toLowerCase();
    var start = currentItem.index + 1;
    if (start === this.treeitems.length) {
      start = 0;
    }
    var index = this.getIndexFirstChars(start, char);
    if (index === -1) {
      index = this.getIndexFirstChars(0, char);
    }
    if (index > -1) {
      this.setFocusToItem(this.treeitems[index]);
    }
  };

  TreeNavController.prototype.findTreeItems = function () {
    var self = this;
    var findItems = function (el, group) {
      var ti = group;
      var curr = el.firstElementChild;
      while (curr) {
        if (curr.tagName === "A" || curr.tagName === "SPAN") {
          ti = new TreeItem(curr, self, group);
          self.treeitems.push(ti);
          self.firstChars.push(ti.label.substring(0, 1).toLowerCase());
        }
        if (curr.firstElementChild) {
          findItems(curr, ti);
        }
        curr = curr.nextElementSibling;
      }
    };
    findItems(this.el, null);
    for (var i = 0; i < this.treeitems.length; i++) {
      this.treeitems[i].index = i;
    }
  };

  TreeNavController.prototype.updateVisibleTreeitems = function () {
    this.firstTreeitem = this.treeitems[0];
    for (var i = 0; i < this.treeitems.length; i++) {
      var ti = this.treeitems[i];
      var parent = ti.groupTreeitem;
      ti.isVisible = true;
      while (parent && parent.el !== this.el) {
        if (!parent.isExpanded()) {
          ti.isVisible = false;
        }
        parent = parent.groupTreeitem;
      }
      if (ti.isVisible) {
        this.lastTreeitem = ti;
      }
    }
  };

  TreeNavController.prototype.setFocusToItem = function (treeitem, focusEl) {
    if (focusEl === undefined) {
      focusEl = true;
    }
    treeitem.el.tabIndex = 0;
    if (focusEl) {
      treeitem.el.focus();
    }
    for (var i = 0; i < this.treeitems.length; i++) {
      if (this.treeitems[i] !== treeitem) {
        this.treeitems[i].el.tabIndex = -1;
      }
    }
  };

  TreeNavController.prototype.getIndexFirstChars = function (startIndex, char) {
    for (var i = startIndex; i < this.firstChars.length; i++) {
      if (this.treeitems[i].isVisible && char === this.firstChars[i]) {
        return i;
      }
    }
    return -1;
  };

  var controller = new TreeNavController(tree);

  // -- mobile <select> fallback -------------------------------------------

  // Build the mobile jump menu from the same tree, indenting nested entries so
  // the type/method hierarchy stays readable in a flat <select>. Selecting an
  // option navigates to its anchor, which also lets the browser scroll to it.
  var select = document.querySelector(".js-mainNavMobile select");
  if (select) {
    var placeholder = document.createElement("option");
    placeholder.textContent = "Outline";
    placeholder.value = "";
    select.appendChild(placeholder);

    for (var i = 0; i < controller.treeitems.length; i++) {
      var ti = controller.treeitems[i];
      var href = ti.el.getAttribute("href");
      if (!href || href.charAt(0) !== "#") {
        continue;
      }
      var opt = document.createElement("option");
      var indent = new Array(Math.max(0, ti.depth - 1) + 1).join("  ");
      opt.textContent = indent + ti.label;
      opt.value = href;
      select.appendChild(opt);
    }

    select.addEventListener("change", function () {
      if (select.value) {
        // Navigating to "#id" updates the hash and scrolls, matching a click on
        // the corresponding tree link. The id is decoded so an encoded anchor
        // still resolves.
        window.location.hash = decodeURIComponent(select.value.slice(1));
      }
    });
  }
})();
