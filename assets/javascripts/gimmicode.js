// Nothing really fancy is this quick'n'dirty code.
$(document).ready(function() {

  // Private vars.
	var lastEnteredCharacter = "";
	var lock                 = false;
	var receivedData         = null;
	
	var $unicodeData         = $('.unicode-data');
	var $characterInput      = $('.character-input');
	
	// Inject the received HTML data into the current page.
	function showUnicode() {
		if (receivedData !== null) {
			$unicodeData.html(receivedData);
			
			receivedData = null;
			lock = false;
		}
	}
	
	// Send the inputed character and request the data.
	function postCharacter(event) {
	  var currentCharacter = $characterInput.val()
		
		var shouldPost = !lock && currentCharacter.length > 0 &&
										 currentCharacter != lastEnteredCharacter;
		
	  if (shouldPost) {
	    lock = true;

	    $.get(
	      $(location).attr('href') + 'unicode',
	      {character: currentCharacter},
	      function(data) {
					receivedData = data;
					showUnicode();
					lastEnteredCharacter = currentCharacter;
	      }
	    )
	  }
	}
	
	// Bind to keyup.
	$characterInput.on('keyup', postCharacter);
});