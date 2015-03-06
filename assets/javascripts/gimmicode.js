// Nothing really fancy is this quick'n'dirty code.
//
// A delay is used the first time a request is made; that
// way the HTML only be injected after the animation played.
$(document).ready(function() {

  // Private vars.
	var animationTimeoutId   = null;
		
	var lastEnteredCharacter = "";
	var lock                 = false;
	var receivedData         = null;
	
	var $inputArea           = $('.input-area');
	var $unicode             = $('.unicode');
	var $unicodeData         = $('.unicode-data');
	var $characterInput      = $('.character-input');
	
	// Inject the received HTML data into the current page.
	function showUnicode() {
		$inputArea.removeClass('animate')

		if (receivedData !== null) {
			$unicodeData.html(receivedData);
			$unicode.removeClass('spin')
			
			receivedData = null;
			lock = false;
		}
		
		clearTimeout(animationTimeoutId);
		animationTimeoutId = null;
	}
	
	// Send the inputed character and request the data.
	function postCharacter() {
	  var currentCharacter = $characterInput.val()
		
		var shouldPost = !lock && currentCharacter.length > 0 &&
										 currentCharacter != lastEnteredCharacter;
		
	  if (shouldPost) {
	    lock = true;

			if (!$inputArea.hasClass('packed')) {
				$inputArea.addClass('animate').addClass('packed');
				animationTimeoutId = setTimeout(showUnicode, 1000);
			}

			$unicode.addClass('spin')

	    $.get(
	      $(location).attr('href') + 'unicode',
	      {character: currentCharacter},
	      function(data) {
					receivedData = data;
					if (animationTimeoutId === null) showUnicode(data);
					lastEnteredCharacter = currentCharacter;
	      }
	    )
	  }
	}
	
	// Bind to keyup.
	$characterInput.on('keyup', postCharacter);
});