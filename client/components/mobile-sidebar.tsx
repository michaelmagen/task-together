import Sidebar from "./sidebar"
import {
    Sheet,
    SheetContent,
    SheetTrigger,
} from "@/components/ui/sheet"
import { Button } from "@/components/ui/button"
import { AlignJustify } from "lucide-react"

// TODO: Make sidebar close when a list is selected
export function MobileSidebar() {
    return (
        <Sheet>
            <SheetTrigger asChild>
                <Button variant="outline" size="icon" className="m-2 md:hidden" >
                    <AlignJustify />
                </Button>
            </SheetTrigger>
            <SheetContent side="left" className="flex overflow-hidden">
                <Sidebar />
            </SheetContent>
        </Sheet>
    )
}
