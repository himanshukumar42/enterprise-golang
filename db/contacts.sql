
CREATE TABLE public.contacts (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES public.user(id),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL,
    country_code VARCHAR(10),
    mobile_number VARCHAR(15),
    events_notification VARCHAR(50) CHECK(events_notification IN ('all_app_users', 'groups')),
    groups TEXT[],
    event_types TEXT[] CHECK (event_types <@ ARRAY['sos', 'timer', 'all']),
    status VARCHAR(20) CHECK(status IN ('active', 'inactive'))
    CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
)


CREATE TABLE public.contacts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255),
    country_code VARCHAR(10),
    mobile_number VARCHAR(20),
    events_notification VARCHAR(50),
    groups VARCHAR(255),
    events_type VARCHAR(255),
    stat VARCHAR(50),
    FOREIGN KEY (user_id) REFERENCES public.user(id)
);