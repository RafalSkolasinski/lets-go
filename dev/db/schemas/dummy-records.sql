-- Add some dummy records (which we'll use in the next couple of chapters).
INSERT INTO
    snippets (title, content, created, expires)
VALUES
    (
        'O Snail',
        'O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
    );

INSERT INTO
    snippets (title, content, created, expires)
VALUES
    (
        'Over the wintry forest',
        'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n- Natsume Soseki',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
    );

INSERT INTO
    snippets (title, content, created, expires)
VALUES
    (
        'First autumn morning',
        'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n- Murakami Kijo',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
    );
