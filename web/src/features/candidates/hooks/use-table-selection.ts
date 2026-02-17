import { useState, useCallback } from "react";

export function useTableSelection<T extends { id: string }>(items: T[]) {
    const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());

    const toggleSelection = useCallback((id: string) => {
        const newSelection = new Set(selectedIds);
        if (newSelection.has(id)) {
            newSelection.delete(id);
        } else {
            newSelection.add(id);
        }
        setSelectedIds(newSelection);
    }, [selectedIds]);

    const toggleAll = useCallback(() => {
        if (selectedIds.size === items.length) {
            setSelectedIds(new Set());
        } else {
            setSelectedIds(new Set(items.map(item => item.id)));
        }
    }, [selectedIds.size, items]);

    const clearSelection = useCallback(() => {
        setSelectedIds(new Set());
    }, []);

    const isAllSelected = selectedIds.size === items.length && items.length > 0;

    const selectedItems = items.filter(item => selectedIds.has(item.id));

    return {
        selectedIds,
        setSelectedIds,
        toggleSelection,
        toggleAll,
        clearSelection,
        isAllSelected,
        selectedItems,
    };
}
