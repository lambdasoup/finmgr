export function connect(app) {
  app.ports.getPushState.subscribe(function() {
    var value;
    if (!available()) {
      value = "NotAvailable";
    } else {
      value = "Available";
    }

    app.ports.setPushState.send(value);
  });
}

function available() {
  if (!('serviceWorker' in navigator)) {
    return false;
  }
  if (!('PushManager' in window)) {
    return false;
  }
  return true;
}
