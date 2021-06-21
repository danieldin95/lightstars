


function updateText() {
	let i18n = $.i18n();
	i18n.locale = navigator.language || navigator.userLanguage;
	console.log(i18n.locale);
	i18n.load( 'i18n/my/' + i18n.locale + '.json', i18n.locale )
		.done(function () {
			$( '.result' ).text( $.i18n("hi") );
		} );
}

// Enable debug
$.i18n.debug = true;

$( document ).ready( function ( $ ) {
	'use strict';
	updateText();
} );
