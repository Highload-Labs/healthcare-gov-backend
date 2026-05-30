CREATE TABLE IF NOT EXISTS coverages
(
    id
    UUID
    PRIMARY
    KEY
    DEFAULT
    gen_random_uuid
(
),
    state varchar
(
    60
) NOT NULL,
    zipcode_start char
(
    5
) NOT NULL,
    zipcode_end char
(
    5
) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             )