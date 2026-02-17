"use client";

import React, { Component, ErrorInfo, ReactNode } from "react";
import { AlertCircle, RefreshCw } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { useTranslations } from "next-intl";

interface Props {
    children: ReactNode;
    fallback?: ReactNode;
}

interface State {
    hasError: boolean;
    error: Error | null;
}

// Wrapper to use hooks (useTranslations) with Class Component
export function WidgetErrorBoundary(props: Props) {
    const t = useTranslations("Errors.Widget");
    return <ErrorBoundaryInner {...props} t={t} />;
}

class ErrorBoundaryInner extends Component<Props & { t: any }, State> {
    public state: State = {
        hasError: false,
        error: null,
    };

    public static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error };
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error("Uncaught error:", error, errorInfo);
    }

    public render() {
        if (this.state.hasError) {
            if (this.props.fallback) {
                return this.props.fallback;
            }

            const { t } = this.props;

            return (
                <Card className="border-destructive/50 bg-destructive/5 h-full min-h-[200px] flex items-center justify-center">
                    <CardContent className="flex flex-col items-center gap-4 p-6 text-center">
                        <AlertCircle className="h-8 w-8 text-destructive" />
                        <div className="space-y-1">
                            <h3 className="font-semibold text-destructive">{t("title")}</h3>
                            <p className="text-sm text-destructive/80 max-w-[200px]">
                                {this.state.error?.message || "Unknown error"}
                            </p>
                        </div>
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => this.setState({ hasError: false })}
                            className="text-destructive border-destructive/30 hover:bg-destructive/10"
                        >
                            <RefreshCw className="mr-2 h-3 w-3" />
                            {t("retry")}
                        </Button>
                    </CardContent>
                </Card>
            );
        }

        return this.props.children;
    }
}
