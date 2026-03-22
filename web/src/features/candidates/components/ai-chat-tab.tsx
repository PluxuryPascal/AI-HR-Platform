import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Send } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState, useRef, useEffect } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { useCandidateChat, useGetChatSessions, useGetChatHistory, ChatMessage, ChatSession } from "../api/use-candidate-ai";
import { useParams } from "next/navigation";
import { toast } from "sonner";

type Message = {
    id: string;
    text: string;
    sender: "user" | "ai";
    timestamp: Date;
};

export const AIChatTab = () => {
    const t = useTranslations("CandidateProfile");
    const params = useParams();
    const candidateId = params.candidateId as string;

    const [messages, setMessages] = useState<Message[]>([]);
    const [inputValue, setInputValue] = useState("");
    const { mutate: sendMessage, isPending: isTyping } = useCandidateChat();
    const { data: sessions } = useGetChatSessions();

    const candidateSession = sessions?.find((s: ChatSession) => s.target_candidate_id === candidateId);
    const { data: history, isLoading: isLoadingHistory } = useGetChatHistory(candidateSession?.id || "");

    useEffect(() => {
        if (history && history.length > 0) {
            setMessages(history.map((m: ChatMessage) => ({
                id: m.id || Math.random().toString(),
                text: m.content,
                sender: m.role === "user" ? "user" : "ai",
                timestamp: m.created_at ? new Date(m.created_at) : new Date(),
            })));
        } else if (!isLoadingHistory) {
            setMessages([
                {
                    id: "welcome",
                    text: "Привет! Я проанализировал профиль этого кандидата. Спрашивайте меня о чем угодно касательно его опыта или навыков.",
                    sender: "ai",
                    timestamp: new Date(),
                },
            ]);
        }
    }, [history, isLoadingHistory]);
    const scrollRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (scrollRef.current) {
            scrollRef.current.scrollIntoView({ behavior: "smooth" });
        }
    }, [messages, isTyping]);

    const handleSend = () => {
        if (!inputValue.trim() || isTyping) return;

        const userText = inputValue;
        const newUserMessage: Message = {
            id: Date.now().toString(),
            text: userText,
            sender: "user",
            timestamp: new Date(),
        };

        setMessages((prev) => [...prev, newUserMessage]);
        setInputValue("");

        sendMessage({
            candidate_id: candidateId,
            question: userText,
            locale: "ru"
        }, {
            onSuccess: (data) => {
                const newAiMessage: Message = {
                    id: Date.now().toString(),
                    text: data.answer,
                    sender: "ai",
                    timestamp: new Date(),
                };
                setMessages((prev) => [...prev, newAiMessage]);
            },
            onError: () => {
                toast.error("Ошибка при получении ответа от ИИ");
            }
        });
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
