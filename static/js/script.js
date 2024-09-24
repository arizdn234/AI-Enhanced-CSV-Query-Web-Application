// smooth scroll
const sections = document.querySelectorAll('section');
const navLinks = document.querySelectorAll('a[href^="#"]');
console.log(navLinks);


window.addEventListener('scroll', () => {
    let current = '';

    sections.forEach(section => {
        const sectionTop = section.offsetTop;
        const sectionHeight = section.clientHeight;

        if (pageYOffset >= sectionTop - sectionHeight / 3) {
            current = section.getAttribute('id');
        }
    });

    navLinks.forEach(link => {
        link.classList.remove('active');
        if (link.getAttribute('href') === `#${current}`) {
            link.classList.add('active');
        }
    });
});

document
    .getElementById("question-form")
    .addEventListener("submit", function (event) {
        event.preventDefault();
        const question = document.getElementById("question").value;

        fetch("/ask", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ question }),
        })
            .then((response) => response.json())
            .then((data) => {
                const responseSection =
                    document.getElementById("response-section");

                const card = document.createElement("div");
                card.className = "card";

                const answerHeading = document.createElement("h3");
                answerHeading.textContent = "Answer";
                card.appendChild(answerHeading);

                const answerPara = document.createElement("p");
                const friendly_answer = `Answer ${data.answer}\t ,Aggregator ${data.aggregator}`;
                answerPara.textContent = friendly_answer;
                card.appendChild(answerPara);

                responseSection.appendChild(card);
            })
            .catch((error) => console.error("Error:", error));
    });

document
    .getElementById("recommend-form")
    .addEventListener("submit", function (event) {
        event.preventDefault();
        const text = document.getElementById("text").value;

        fetch("/recommend", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ text }),
        })
            .then((response) => response.json())
            .then((data) => {
                const recommendationsList = document.getElementById(
                    "recommendations-list"
                );

                const recommendationDiv = document.createElement("div");
                recommendationDiv.classList.add("card");

                const inputPara = document.createElement("p");
                inputPara.classList.add("recommendation-input");
                inputPara.innerHTML = `<strong>Input:</strong> ${data.input}`;
                recommendationDiv.appendChild(inputPara);

                const recommendationPara = document.createElement("div");
                recommendationPara.classList.add("recommendation-output");
                const recommendationText = data.recommendation.Parts.join(" ");
                recommendationPara.innerHTML =
                    convertMarkdownToHTML(recommendationText);
                recommendationDiv.appendChild(recommendationPara);

                recommendationsList.appendChild(recommendationDiv);
            })
            .catch((error) => console.error("Error:", error));
    });

function convertMarkdownToHTML(markdown) {
    const rules = [
        { regex: /###### (.*$)/gim, replacement: "<h6>$1</h6>" },
        { regex: /##### (.*$)/gim, replacement: "<h5>$1</h5>" },
        { regex: /#### (.*$)/gim, replacement: "<h4>$1</h4>" },
        { regex: /### (.*$)/gim, replacement: "<h3>$1</h3>" },
        { regex: /## (.*$)/gim, replacement: "<h2>$1</h2>" },
        { regex: /# (.*$)/gim, replacement: "<h1>$1</h1>" },
        { regex: /\*\*(.*)\*\*/gim, replacement: "<strong>$1</strong>" },
        { regex: /\*(.*)\*/gim, replacement: "<em>$1</em>" },
        { regex: /\n$/gim, replacement: "<br />" },
        { regex: /^\> (.*$)/gim, replacement: "<blockquote>$1</blockquote>" },
        { regex: /\n\*(.*)/gim, replacement: "<ul>\n<li>$1</li>\n</ul>" },
        {
            regex: /\n[0-9]+\.(.*)/gim,
            replacement: "<ol>\n<li>$1</li>\n</ol>",
        },
    ];
    let html = markdown;
    rules.forEach((rule) => {
        html = html.replace(rule.regex, rule.replacement);
    });
    return html.trim();
}

function logout() {
    window.location.href = "/logout";
}