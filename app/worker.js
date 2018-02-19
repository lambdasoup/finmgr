self.addEventListener('push', function(event) {
  console.log('event:', event);
  console.log('data:', event.data.text());
});
