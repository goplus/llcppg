// Outline sidebar behaviour for `llcppg -doc`.
//
// internal/godoc renders the documentation outline as a nested <ul> tree using
// pkg.go.dev's class names and #anchors. This script enhances that static tree
// with the two behaviours a reader expects from pkg.go.dev's sidebar, without
// any build step or network access:
//
//   1. Scroll-spy: the outline link whose target section is currently in view
//      is marked current (aria-current="true"), so the sidebar tracks the page
//      as it scrolls.
//   2. Keyboard navigation: arrow keys move between visible tree items and
//      Enter/Space follows the focused link, matching an ARIA tree widget.
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

  // links are every outline anchor in document order. treeItems are the ones
  // that participate in keyboard navigation (the same set; separated only for
  // clarity).
  var links = Array.prototype.slice.call(tree.querySelectorAll("a[href^='#']"));
  if (links.length === 0) {
    return;
  }

  // -- scroll-spy ---------------------------------------------------------

  // targets maps a section element to the outline link that points at it, so an
  // IntersectionObserver callback can mark the right link current. Only anchors
  // that resolve to an element on the page are tracked; "Source Files" and
  // "Documentation" section wrappers resolve like any other.
  var byId = {};
  links.forEach(function (a) {
    var id = decodeURIComponent(a.getAttribute("href").slice(1));
    var el = document.getElementById(id);
    if (el) {
      byId[id] = { link: a, el: el };
    }
  });

  var ids = Object.keys(byId);
  if (ids.length > 0 && "IntersectionObserver" in window) {
    // visible tracks which tracked sections are currently intersecting; the
    // topmost visible one wins so the sidebar highlights the section the reader
    // is looking at rather than one scrolled past.
    var visible = {};
    var setCurrent = function () {
      var bestId = null;
      var bestTop = Infinity;
      ids.forEach(function (id) {
        if (visible[id]) {
          var top = byId[id].el.getBoundingClientRect().top;
          if (top < bestTop) {
            bestTop = top;
            bestId = id;
          }
        }
      });
      links.forEach(function (a) {
        a.removeAttribute("aria-current");
      });
      if (bestId) {
        byId[bestId].link.setAttribute("aria-current", "true");
      }
    };

    var observer = new IntersectionObserver(
      function (entries) {
        entries.forEach(function (entry) {
          var id = entry.target.id;
          if (entry.isIntersecting) {
            visible[id] = true;
          } else {
            delete visible[id];
          }
        });
        setCurrent();
      },
      // A top margin keeps a section "current" while its heading sits just
      // under the top of the viewport, which reads more naturally than
      // switching exactly at the fold.
      { rootMargin: "0px 0px -70% 0px", threshold: 0 }
    );
    ids.forEach(function (id) {
      observer.observe(byId[id].el);
    });
  }

  // -- keyboard navigation (ARIA tree) ------------------------------------

  // Make the tree focusable as a single widget: the first item is in the tab
  // order, the rest are reachable with the arrow keys (roving tabindex).
  links.forEach(function (a, i) {
    a.setAttribute("tabindex", i === 0 ? "0" : "-1");
  });

  var focusAt = function (i) {
    if (i < 0) {
      i = 0;
    }
    if (i >= links.length) {
      i = links.length - 1;
    }
    links.forEach(function (a, j) {
      a.setAttribute("tabindex", j === i ? "0" : "-1");
    });
    links[i].focus();
  };

  tree.addEventListener("keydown", function (e) {
    var i = links.indexOf(document.activeElement);
    if (i < 0) {
      return;
    }
    switch (e.key) {
      case "ArrowDown":
        e.preventDefault();
        focusAt(i + 1);
        break;
      case "ArrowUp":
        e.preventDefault();
        focusAt(i - 1);
        break;
      case "Home":
        e.preventDefault();
        focusAt(0);
        break;
      case "End":
        e.preventDefault();
        focusAt(links.length - 1);
        break;
      case "Enter":
      case " ":
        e.preventDefault();
        links[i].click();
        break;
      default:
        break;
    }
  });

  // -- mobile <select> fallback -------------------------------------------

  // Build the mobile jump menu from the same tree, indenting nested entries so
  // the type/method hierarchy stays readable in a flat <select>. Selecting an
  // option navigates to its anchor, which also lets the browser scroll to it.
  var select = document.querySelector(".js-mainNavMobile select");
  if (select) {
    var depthOf = function (a) {
      // Depth is the number of <ul> ancestors between the link's <li> and the
      // tree root; it drives the indentation of the option label.
      var depth = 0;
      var node = a.parentNode;
      while (node && node !== tree) {
        if (node.tagName === "UL") {
          depth++;
        }
        node = node.parentNode;
      }
      return Math.max(0, depth - 1);
    };

    var placeholder = document.createElement("option");
    placeholder.textContent = "Outline";
    placeholder.value = "";
    select.appendChild(placeholder);

    links.forEach(function (a) {
      var opt = document.createElement("option");
      var indent = new Array(depthOf(a) + 1).join("  ");
      opt.textContent = indent + a.textContent.trim();
      opt.value = a.getAttribute("href");
      select.appendChild(opt);
    });

    select.addEventListener("change", function () {
      if (select.value) {
        // Navigating to "#id" updates the hash and scrolls, matching a click on
        // the corresponding tree link. The id is decoded so it keys the same way
        // the scroll-spy map does, in case godoc ever emits an encoded anchor.
        window.location.hash = decodeURIComponent(select.value.slice(1));
      }
    });
  }
})();
