"use client";

import React, { useEffect, useRef } from "react";
import { useTheme } from "next-themes";

interface Dot {
    x: number;
    y: number;
    originX: number;
    originY: number;
    vx: number;
    vy: number;
    size: number;
    color: string;
}

export const GlobalDotGridBg = () => {
    const canvasRef = useRef<HTMLCanvasElement>(null);
    const { resolvedTheme } = useTheme();

    useEffect(() => {
        const canvas = canvasRef.current;
        if (!canvas) return;

        const ctx = canvas.getContext("2d");
        if (!ctx) return;

        let animationFrameId: number;
        let dots: Dot[] = [];
        const mouse = { x: -1000, y: -1000 };

        // Physics parameters
        const SPACING = 25;
        const RADIUS = 150; // Interaction radius
        const HOVER_RADIUS = 300; // Visibility radius ("flashlight")
        const FORCE_FACTOR = 1.2;
        const SPRING_FACTOR = 0.05;
        const DAMPING = 0.90;
        const DOT_SIZE = 1.5;

        // determine theme color
        const isDark = resolvedTheme === "dark";
        const dotColor = isDark ? "rgba(255, 255, 255, " : "rgba(0, 0, 0, ";

        const initDots = () => {
            dots = [];
            const cols = Math.ceil(window.innerWidth / SPACING);
            const rows = Math.ceil(window.innerHeight / SPACING);

            for (let i = 0; i < cols; i++) {
                for (let j = 0; j < rows; j++) {
                    const x = i * SPACING;
                    const y = j * SPACING;
                    dots.push({
                        x,
                        y,
                        originX: x,
                        originY: y,
                        vx: 0,
                        vy: 0,
                        size: DOT_SIZE,
                        color: dotColor,
                    });
                }
            }
        };

        const handleResize = () => {
            canvas.width = window.innerWidth;
            canvas.height = window.innerHeight;
            initDots();
        };

        const handleMouseMove = (e: MouseEvent) => {
            mouse.x = e.clientX;
            mouse.y = e.clientY;
        };

        const render = (time: number) => {
            ctx.clearRect(0, 0, canvas.width, canvas.height);

            // Idle movement factor (slow drift)
            const timeFactor = time * 0.0005;

            dots.forEach((dot, index) => {
                // 1. Calculate distance from mouse
                const dx = mouse.x - dot.x;
                const dy = mouse.y - dot.y;
                const dist = Math.sqrt(dx * dx + dy * dy);

                // Visibility Check ("Flashlight" Effect)
                // If too far, don't draw (opacity = 0)
                let alpha = 0;
                if (dist < HOVER_RADIUS) {
                    // Mapping 0-HOVER_RADIUS distance to 0.4-0 opacity
                    // The closer, the more visible
                    alpha = (1 - dist / HOVER_RADIUS) * 0.4;
                }

                // Optimization: Skip physics and drawing if totally invisible
                // BUT, we still want physics if we want them to "spring" into view when we move fast
                // For now, let's always run physics to keep state consistent, but only draw if alpha > 0.01

                // 2. Apply Interaction Force (Repulsion)
                if (dist < RADIUS) {
                    const angle = Math.atan2(dy, dx);
                    // Calculate force: stronger when closer
                    const force = (RADIUS - dist) / RADIUS;
                    const push = force * FORCE_FACTOR;

                    dot.vx -= Math.cos(angle) * push;
                    dot.vy -= Math.sin(angle) * push;
                }

                // 3. Apply Spring Force (Return to Origin)
                // Add minimal idle noise to origin
                const noiseX = Math.sin(timeFactor + index) * 2;
                const noiseY = Math.cos(timeFactor + index) * 2;

                const targetX = dot.originX + noiseX;
                const targetY = dot.originY + noiseY;

                const springDx = targetX - dot.x;
                const springDy = targetY - dot.y;

                dot.vx += springDx * SPRING_FACTOR;
                dot.vy += springDy * SPRING_FACTOR;

                // 4. Apply Damping (Friction)
                dot.vx *= DAMPING;
                dot.vy *= DAMPING;

                // 5. Update Position
                dot.x += dot.vx;
                dot.y += dot.vy;

                // 6. Draw Dot if visible
                if (alpha > 0.01) {
                    ctx.fillStyle = `${dot.color}${alpha})`;
                    ctx.beginPath();
                    ctx.arc(dot.x, dot.y, dot.size, 0, Math.PI * 2);
                    ctx.fill();
                }
            });

            animationFrameId = requestAnimationFrame(render);
        };

        // Initialize
        handleResize();
        window.addEventListener("resize", handleResize);
        window.addEventListener("mousemove", handleMouseMove);
        animationFrameId = requestAnimationFrame(render);

        return () => {
            window.removeEventListener("resize", handleResize);
            window.removeEventListener("mousemove", handleMouseMove);
            cancelAnimationFrame(animationFrameId);
        };
    }, [resolvedTheme]); // Re-run when theme changes

    return (
        <div className="fixed inset-0 -z-10 bg-slate-50 dark:bg-background transition-colors duration-300">
            <canvas ref={canvasRef} className="block w-full h-full" />
            {/* Ambient Gradient Ovelay for depth */}
            <div className="absolute inset-0 bg-gradient-to-tr from-blue-50/20 via-transparent to-indigo-50/20 dark:from-blue-950/20 dark:to-indigo-950/20 pointer-events-none" />
        </div>
    );
};
