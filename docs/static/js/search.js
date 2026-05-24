const input = document.querySelector('#search-input');
const results = document.querySelector('#search-results');
let index = [];

async function loadIndex() {
  if (index.length) return;
  const res = await fetch('/index.json');
  index = await res.json();
}

function escapeHtml(value) {
  return value.replace(/[&<>'"]/g, char => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', "'": '&#39;', '"': '&quot;' }[char]));
}

function render(matches) {
  if (!matches.length) {
    results.hidden = true;
    results.innerHTML = '';
    return;
  }

  results.innerHTML = matches.map(page => `<a href="${page.href}"><strong>${escapeHtml(page.title)}</strong><span>${escapeHtml(page.content.slice(0, 120))}...</span></a>`).join('');
  results.hidden = false;
}

if (input && results) {
  input.addEventListener('input', async event => {
    const query = event.target.value.trim().toLowerCase();
    if (query.length < 2) {
      render([]);
      return;
    }

    await loadIndex();
    render(index.filter(page => `${page.title} ${page.content}`.toLowerCase().includes(query)).slice(0, 6));
  });
}
