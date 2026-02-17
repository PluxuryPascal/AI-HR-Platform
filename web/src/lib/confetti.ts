import confetti from "canvas-confetti";

export const fireOfferConfetti = () => {
    const duration = 1000;
    const end = Date.now() + duration;

    // Colors: Emerald-500 (#10b981) and Blue-500 (#3b82f6)
    const colors = ["#10b981", "#3b82f6"];

    (function frame() {
        confetti({
            particleCount: 2,
            angle: 60,
            spread: 55,
            origin: { x: 0 },
            colors: colors,
        });
        confetti({
            particleCount: 2,
            angle: 120,
            spread: 55,
            origin: { x: 1 },
            colors: colors,
        });

        if (Date.now() < end) {
            requestAnimationFrame(frame);
        }
    })();
};
