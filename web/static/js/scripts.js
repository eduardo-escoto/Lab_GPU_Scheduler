document.addEventListener('DOMContentLoaded', function() {
    // Example of an HTMX interaction to load user data
    const loadUserButton = document.getElementById('load-user');
    loadUserButton.addEventListener('click', function() {
        htmx.ajax('GET', '/users', '#user-data');
    });

    // Example of form submission using HTMX
    const userForm = document.getElementById('user-form');
    userForm.addEventListener('submit', function(event) {
        event.preventDefault();
        htmx.ajax('POST', '/users', '#user-data', { 
            values: htmx.serialize(userForm) 
        });
    });
});