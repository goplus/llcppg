(function () {
  "use strict";

  var marker = document.createElement("div");
  marker.id = "llcppg-live";
  document.addEventListener("DOMContentLoaded", function () {
    document.body.appendChild(marker);
  });

  function setStatus(state) {
    marker.dataset.state = state;
    marker.textContent = state === "live" ? "live" : "disconnected";
  }
  setStatus("connecting");

  // serverId remembers the instance id from the first hello. A later hello with
  // a different id means the server restarted while this tab was open, so the
  // page reloads to pick up whatever the new server serves.
  var serverId = null;

  // reload wraps location.reload() so a live-reload navigation can be told apart
  // from the user closing the tab. In auto-exit mode the page beacons /_doc/bye
  // on pagehide; location.reload() fires pagehide too, so without this flag an
  // automatic update would look like a close and start the server's exit timer
  // (see assets/page.html). We mark the flag on window so the pagehide handler,
  // which lives in a separate inline script, can read it and skip the beacon.
  function reload() {
    window.__llcppgReloading = true;
    location.reload();
  }

  var es = new EventSource("/_events");

  es.addEventListener("hello", function (e) {
    if (serverId === null) {
      serverId = e.data;
    } else if (e.data !== serverId) {
      reload();
      return;
    }
    setStatus("live");
  });

  es.addEventListener("reload", function () {
    // A full reload re-runs the server's single render path, so the live view
    // can never drift from a manual refresh. The browser keeps the scroll
    // position and the URL fragment across the reload.
    reload();
  });

  es.onopen = function () {
    setStatus("live");
  };

  es.onerror = function () {
    // EventSource reconnects on its own; reflect the gap so a stale page is not
    // mistaken for a fresh one.
    setStatus("disconnected");
  };
})();
