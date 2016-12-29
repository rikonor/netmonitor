function init() {
  const filter = { urls: [ "<all_urls>" ] };
  chrome.webRequest.onBeforeRequest.addListener(onBefore, filter);
}

function onBefore({ tabId, method, type, url }) {
  getCurrentTab()
		.then(function({ activeTabId, activeTabUrl }) {
      console.log("x", activeTabId, activeTabUrl);
			if (activeTabId === tabId) {
				fetch("http://localhost:8080", {
					method: 'POST',
					body: JSON.stringify({ siteUrl: activeTabUrl, method, type, url })
				})
			}
		});
}

function getCurrentTab() {
	return new Promise(function(resolve) {
    chrome.tabs.query({ active: true }, function([ { id, url } ]) {
			resolve({ activeTabId: id, activeTabUrl: url });
		});
	});
}

init()
