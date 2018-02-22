//# sourceURL=app.js
'use strict';

import elm from "./elm"
import * as grpc from "./grpc";
import push from "./push";

function pwaSupported() {
  if (!('serviceWorker' in navigator)) {
    return false;
  }
  if (!('PushManager' in window)) {
    return false;
  }
  if (!('showNotification' in ServiceWorkerRegistration.prototype)) {
    return false;
  }
  return true;
}

if (pwaSupported()) {
  var app = elm.Main.fullscreen()
  grpc.connect(app);
  push.connect(app);
} else {
  alert("this app needs a proper browser!");
}
