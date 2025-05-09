document.addEventListener('DOMContentLoaded', () => {
  const container = document.getElementById('container');
  const signUpButton = document.getElementById('signUp');
  const signInButton = document.getElementById('signIn');
  const signupForm = document.getElementById('loginForm')
  const sign = document.getElementById('sing');
  const mainPage = document.getElementById('mainPage');


  signUpButton.addEventListener('click', () => {
    container.classList.add('right-panel-active');
    signupForm.style.display='none'

  });

  signInButton.addEventListener('click', () => {
    container.classList.remove('right-panel-active');
    signupForm.style.display='block'
  });
  sign.addEventListener('click', (e) => {
    e.preventDefault(); 
    container.style.display = 'none'; 
    mainPage.style.display = 'block'; 
  });
  
    
});
