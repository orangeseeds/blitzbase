CREATE TABLE "_base_collection_posts" ("content" TEXT DEFAULT ""  NOT NULL, "created" TIMESTAMP DEFAULT CURRENT_TIMESTAMP  NOT NULL, "id" INTEGER PRIMARY KEY, "title" TEXT DEFAULT ""  NOT NULL, "updated" TIMESTAMP DEFAULT CURRENT_TIMESTAMP  NOT NULL, "user" INTEGER NOT NULL  REFERENCES users(id), "views_count" INTEGER DEFAULT 0  NOT NULL);