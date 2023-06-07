document.body.onload = function() {
	setTimeout(function(){
		let preloader = document.getElementById('page-preloader');
		if( !preloader.classList.contains('done'))
		{
			preloader.classList.add('done')
		}
	}, 1000);
}

let menuBtn = document.querySelector('.menu-btn');
let menu = document.querySelector('.menu');
let overFlow = document.querySelector('body');
menuBtn.addEventListener('click', function(){
	menuBtn.classList.toggle('active');
	menu.classList.toggle('active');
	overFlow.classList.toggle('noscroll');
})

