bun build ./Scripts/Main.ts --outdir ./Out/Scripts --watch --minify &
sass --watch Styles/Main.scss:Out/Styles/Main.min.css --style=compressed --silence-deprecation=import
wait
