        CREATE TABLE IF NOT EXISTS basestats(
            name TEXT,
            class TEXT,
            level INTEGER,
            hp INTEGER,
            str INTEGER,
            mag INTEGER,
            skl INTEGER,
            spd INTEGER,
            lck INTEGER,
            def INTEGER,
            res INTEGER,
            mov INTEGER
        );

        CREATE TABLE IF NOT EXISTS basegrowths(
            name TEXT,
            hp INTEGER,
            str INTEGER,
            mag INTEGER,
            skl INTEGER,
            spd INTEGER,
            lck INTEGER,
            def INTEGER,
            res INTEGER
        );

        CREATE TABLE IF NOT EXISTS asset(
            asset TEXT,
            hp TEXT,
            str TEXT,
            mag TEXT,
            skl TEXT,
            spd TEXT,
            lck TEXT,
            def TEXT,
            res TEXT,
            affinity INT,
            growth INT
        );

        CREATE TABLE IF NOT EXISTS classes(
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT
        );

        CREATE TABLE IF NOT EXISTS classbase(
            class TEXT,
            hp INTEGER,
            str INTEGER,
            mag INTEGER,
            skl INTEGER,
            spd INTEGER,
            def INTEGER,
            res INTEGER,
            mov INTEGER,
            rank TEXT
        );

        CREATE TABLE IF NOT EXISTS skills(
            skill TEXT,
            effect TEXT,
            activation TEXT,
            class TEXT,
            level INTEGER
        );

        CREATE TABLE IF NOT EXISTS character_assets(
            name TEXT,
            str INTEGER,
            mag INTEGER,
            skl INTEGER,
            spd INTEGER,
            lck INTEGER,
            def INTEGER,
            res INTEGER
        );
