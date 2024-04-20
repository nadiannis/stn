const linkDeleteButton = document.querySelector('.btn-link-delete');
const id = location.pathname.split('/')[2];

linkDeleteButton.addEventListener('click', () => deleteLink(id));
