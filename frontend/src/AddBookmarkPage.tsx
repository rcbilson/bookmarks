import React, { useState } from "react";
import { Input, Button } from '@chakra-ui/react';
import { Toaster, toaster } from "@/components/ui/toaster"
import axios from "axios";
import { useQueryClient } from '@tanstack/react-query';

const AddBookmarkPage: React.FC = () => {
    const [url, setUrl] = useState("");
    const queryClient = useQueryClient();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            new URL(url);
        } catch (e) {
            toaster.create({
                title: "Invalid URL",
                type: "error",
              });
            return;
        }
        const promise = axios.post("/api/add?url=" + encodeURIComponent(url)).then(
            (_) => {
                queryClient.invalidateQueries({ queryKey: ['bookmarkList'] });
                setUrl("");
            }
        );
        toaster.promise(promise, {
            success: {
                title: "Bookmark Added!",
                description: "Looks great",
            },
            error: {
                title: "Add failed",
                description: "Something went wrong",
            },
            loading: { title: "Adding...", description: "Please wait" },
        })
    };

    return (
        <>
            <Toaster />
            <form onSubmit={handleSubmit}>
                <Input
                    value={url}
                    onChange={(e) => setUrl(e.target.value)}
                    placeholder="Enter URL to bookmark"
                    mb={4}
                />
                <Button type="submit">Add Bookmark</Button>
            </form>
        </>
    );
};

export default AddBookmarkPage;
