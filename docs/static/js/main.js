const navToggle = document.querySelector('.nav-toggle');
const sidebar = document.querySelector('.sidebar');

if (navToggle && sidebar) {
  navToggle.addEventListener('click', () => {
    const open = sidebar.classList.toggle('open');
    navToggle.setAttribute('aria-expanded', String(open));
  });
}
