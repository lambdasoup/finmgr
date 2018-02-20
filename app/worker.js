self.addEventListener('push', function(event) {
  console.log('data:', event.data.text());
});
