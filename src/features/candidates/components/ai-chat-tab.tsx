import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Send } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState, useRef, useEffect } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

type Message = {
    id: string;
    text: string;
    sender: "user" | "ai";
    timestamp: Date;
};

export const AIChatTab = () => {
    const t = useTranslations("CandidateProfile");
    const [messages, setMessages] = useState<Message[]>([
        {
            id: "1",
            text: "Hello! I've analyzed this candidate's profile. Ask me anything about their experience or skills.",
            sender: "ai",
            timestamp: new Date(),
        },
    ]);
    const [inputValue, setInputValue] = useState("");
    const [isTyping, setIsTyping] = useState(false);
    const scrollRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollIntoView({ behavior: "smooth" });
        }
    }, [messages, isTyping]);

    const handleSend = () => {
        if (!inputValue.trim()) return;

        const newUserMessage: Message = {
            id: Date.now().toString(),
            text: inputValue,
            sender: "user",
            timestamp: new Date(),
        };

        setMessages((prev) => [...prev, newUserMessage]);
        setInputValue("");
        setIsTyping(true);

        // Simulate AI response
        setTimeout(() => {
            const responses = [
                "Based on the resume, the candidate has significant experience in that area.",
                "Yes, they mention that skill in their recent project at TechCorp.",
                "That's a good question. While not explicitly stated, their background suggests familiarity with it.",
                "The candidate's strength lies more in frontend development than backend engineering.",
            ];
            const randomResponse =
                responses[Math.floor(Math.random() * responses.length)];

            const newAiMessage: Message = {
                id: (Date.now() + 1).toString(),
                text: randomResponse,
                sender: "ai",
                timestamp: new Date(),
            };

            setMessages((prev) => [...prev, newAiMessage]);
            setIsTyping(false);
        }, 1500);
    };

    return (
        <div className="flex flex-col h-full bg-gray-50/50 dark:bg-muted/10">
            <ScrollArea className="flex-1 min-h-0">
                <div className="space-y-4 max-w-3xl mx-auto px-6 py-4">
                    {messages.map((msg) => (
                        <div
                            key={msg.id}
                            className={`flex items-start gap-3 ${msg.sender === "user" ? "flex-row-reverse" : "flex-row"
                                }`}
                        >
                            <Avatar className="w-8 h-8 mt-1 border">
                                <AvatarFallback className={msg.sender === "ai" ? "bg-blue-100 text-blue-600 dark:bg-blue-900 dark:text-blue-200" : "bg-gray-100 dark:bg-gray-800"}>
                                    {msg.sender === "ai" ? "AI" : "ME"}
                                </AvatarFallback>
                            </Avatar>
                            <div
                                className={`p-3 rounded-lg max-w-[80%] text-sm ${msg.sender === "user"
                                    ? "bg-blue-600 text-white rounded-tr-none"
                                    : "bg-white dark:bg-card border dark:border-border rounded-tl-none shadow-sm text-gray-800 dark:text-gray-100"
                                    }`}
                            >
                                {msg.text}
                            </div>
                        </div>
                    ))}
                    {isTyping && (
                        <div className="flex items-center gap-3">
                            <Avatar className="w-8 h-8 mt-1 border">
                                <AvatarFallback className="bg-blue-100 text-blue-600 dark:bg-blue-900 dark:text-blue-200">AI</AvatarFallback>
                            </Avatar>
                            <div className="bg-white dark:bg-card border dark:border-border p-3 rounded-lg rounded-tl-none shadow-sm">
                                <div className="flex space-x-1">
                                    <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: "0ms" }} />
                                    <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: "150ms" }} />
                                    <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: "300ms" }} />
                                </div>
                            </div>
                        </div>
                    )}
                    <div ref={scrollRef} />
                </div>
            </ScrollArea>
            <div className="p-4 bg-white dark:bg-card border-t dark:border-border">
                <form
                    onSubmit={(e) => {
                        e.preventDefault();
                        handleSend();
                    }}
                    className="flex gap-2 max-w-3xl mx-auto"
                >
                    <Input
                        value={inputValue}
                        onChange={(e) => setInputValue(e.target.value)}
                        placeholder={t("chatPlaceholder")}
                        className="flex-1 bg-background dark:bg-background"
                    />
                    <Button type="submit" size="icon" disabled={!inputValue.trim() || isTyping}>
                        <Send className="w-4 h-4" />
                    </Button>
                </form>
            </div>
        </div>
    );
};
