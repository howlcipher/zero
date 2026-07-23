document.addEventListener('DOMContentLoaded', () => {
    // Theme toggle logic
    const themeToggleBtn = document.getElementById('theme-toggle');
    const body = document.body;

    // Check localStorage for theme preference
    // Default is dark mode, so we only need to add 'light-mode' class if saved as light
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'light') {
        body.classList.add('light-mode');
    }

    themeToggleBtn.addEventListener('click', () => {
        body.classList.toggle('light-mode');
        if (body.classList.contains('light-mode')) {
            localStorage.setItem('theme', 'light');
        } else {
            localStorage.setItem('theme', 'dark');
        }
    });

    // Retro text decoding/scramble effect on main title
    const mainTitle = document.getElementById('main-title');
    const originalText = mainTitle.innerText;
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%^&*()_+';
    let iterations = 0;

    const interval = setInterval(() => {
        mainTitle.innerText = originalText
            .split('')
            .map((letter, index) => {
                // If it's a space, keep it a space
                if (letter === ' ') return ' ';
                
                if (index < iterations) {
                    return originalText[index];
                }
                return chars[Math.floor(Math.random() * chars.length)];
            })
            .join('');

        if (iterations >= originalText.length) {
            clearInterval(interval);
        }

        iterations += 1 / 3;
    }, 30);
});
