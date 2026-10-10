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

  var es = new EventSource("/_events");

  es.addEventListener("hello", function (e) {
    if (serverId === null) {
      serverId = e.data;
    } else if (e.data !== serverId) {
      location.reload();
      return;
    }
    setStatus("live");
  });

  es.addEventListener("reload", function () {
    // A full reload re-runs the server's single render path, so the live view
    // can never drift from a manual refresh. The browser keeps the scroll
    // position and the URL fragment across the reload.
    location.reload();
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
