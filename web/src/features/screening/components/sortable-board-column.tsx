import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { BoardColumn } from "./board-column";
import { CandidateCard } from "../types";

interface SortableBoardColumnProps {
    id: string;
    title: string;
    candidates: CandidateCard[];
    selectedCandidateIds: string[];
    onToggleSelection: (id: string) => void;
    isSelectionMode: boolean;
    isEditMode: boolean;
    onEdit: () => void;
    onDelete: () => void;
}

export function SortableBoardColumn(props: SortableBoardColumnProps) {
    const {
        attributes,
        listeners,
        setNodeRef,
        transform,
        transition,
        isDragging,
    } = useSortable({ 
        id: props.id,
        disabled: !props.isEditMode
    });

    const style = {
        transform: CSS.Translate.toString(transform),
        transition,
        opacity: isDragging ? 0.5 : 1,
    };

    return (
        <div ref={setNodeRef} style={style} {...attributes} {...(props.isEditMode ? listeners : {})}>
            <BoardColumn {...props} />
        </div>
    );
}
