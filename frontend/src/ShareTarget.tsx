import { useEffect, useRef, useState } from "react";
import axios from "axios";
import { Toaster, toaster } from "@/components/ui/toaster"

function isValidUrl(s: string): boolean {
  try {
    const u = new URL(s);
    return u.protocol === "http:" || u.protocol === "https:";
  } catch {
    return false;
  }
}

interface ShareTargetProps {
  onDone: () => void;
}

// This component acts as a PWA share target. The manifest's share_target
// uses GET, so the shared data arrives as query parameters.
export default function ShareTarget({ onDone }: ShareTargetProps) {
  const params = new URLSearchParams(window.location.search);
  // Most Android apps put the link in "text"; fall back to "url" for
  // those that populate it instead.
  const text = params.get("text")?.trim() || params.get("url")?.trim() || "";
  const title = params.get("title")?.trim() || "";

  const [error, setError] = useState<string | null>(
    !text
      ? "No link was shared."
      : !isValidUrl(text)
        ? `The shared text is not a valid URL: "${text}"`
        : null
  );

  // Guard against StrictMode's double-invoked effects adding twice.
  const started = useRef(false);

  useEffect(() => {
    if (error || started.current) return;
    started.current = true;

    let query = "url=" + encodeURIComponent(text);
    if (title) {
      query += "&title=" + encodeURIComponent(title);
    }
    axios.post("/api/add?" + query).then(
      () => {
        toaster.create({ title: "Bookmark added!", type: "success" });
        onDone();
      },
      (e) => {
        setError(`Add failed: ${e instanceof Error ? e.message : String(e)}`);
      }
    );
  }, [error, text, title, onDone]);

  if (error) {
    return (
      <>
        <Toaster />
        <div style={{ padding: "2em", textAlign: "center" }}>
          <h2>Could not add bookmark</h2>
          <p style={{ color: '#d32f2f' }}>{error}</p>
          <p>Please share a valid web link (starting with http:// or https://).</p>
        </div>
      </>
    );
  }

  return (
    <>
      <Toaster />
      <div style={{ padding: "2em", textAlign: "center" }}>
        <h2>Adding shared link...</h2>
        <p>The url received was {window.location.href}</p>
      </div>
    </>
  );
}
